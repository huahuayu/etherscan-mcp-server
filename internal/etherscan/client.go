package etherscan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents an Etherscan API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Response is the standard response format from Etherscan API
type Response struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// JSONRPCResponse is the format for JSON-RPC responses
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents an error in a JSON-RPC response
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrNotFreeAPI is returned when Etherscan API denies access for non-free chains
var ErrNotFreeAPI = errors.New("etherscan API: this chain requires a paid plan")

// IsNotFreeAPIError checks if an error is caused by a non-free API response
func IsNotFreeAPIError(err error) bool {
	return errors.Is(err, ErrNotFreeAPI)
}

// Error represents an API error
type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("etherscan API error: %s - %s", e.Status, e.Message)
}

// NewClient creates a new Etherscan client
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.etherscan.io/v2/api",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Request performs a GET request to the Etherscan API
func (c *Client) Request(chainID string, module, action string, params map[string]string) (json.RawMessage, error) {
	// Create URL values
	values := url.Values{}
	values.Set("module", module)
	values.Set("action", action)
	values.Set("apikey", c.apiKey)
	values.Set("chainid", chainID)

	// Add additional parameters
	for k, v := range params {
		values.Set(k, v)
	}

	// Create request URL with chainID as a query parameter
	// Format: https://api.etherscan.io/v2/api?chainid=${chainid}&${other-params}
	requestURL := fmt.Sprintf("%s?%s", c.baseURL, values.Encode())

	// Create request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// First, try to parse as a JSON-RPC response (for proxy module)
	var jsonRPCResponse JSONRPCResponse
	if err := json.Unmarshal(body, &jsonRPCResponse); err == nil && jsonRPCResponse.JSONRPC != "" {
		// This is a JSON-RPC response
		if jsonRPCResponse.Error != nil {
			return nil, fmt.Errorf("JSON-RPC error: %d - %s", jsonRPCResponse.Error.Code, jsonRPCResponse.Error.Message)
		}
		return jsonRPCResponse.Result, nil
	}

	// If not a JSON-RPC response, parse as standard Etherscan API response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in standard response
	if response.Status != "1" && response.Status != "" {
		// Detect non-free API error (NOTOK typically means the chain requires a paid plan)
		if response.Status == "0" && response.Message == "NOTOK" {
			return nil, fmt.Errorf("%w: %s", ErrNotFreeAPI, string(response.Result))
		}
		return nil, &Error{
			Status:  response.Status,
			Message: response.Message,
		}
	}

	return response.Result, nil
}

// GetAccountBalance gets the balance of an account on a specific blockchain
func (c *Client) GetAccountBalance(chainID, address string) (string, error) {
	params := map[string]string{
		"address": address,
		"tag":     "latest",
	}

	result, err := c.Request(chainID, "account", "balance", params)
	if err != nil {
		return "", err
	}

	var balance string
	if err := json.Unmarshal(result, &balance); err != nil {
		return "", fmt.Errorf("failed to parse balance: %w", err)
	}

	return balance, nil
}

// GetBlockByNumber gets block information by block number
func (c *Client) GetBlockByNumber(chainID, blockNumber string) (json.RawMessage, error) {
	// For non-proxy API, we don't need to convert to hex format
	params := map[string]string{
		"blockno": blockNumber,
	}

	return c.Request(chainID, "block", "getblockreward", params)
}

// GetBlockByNumberRaw gets block information by block number using the raw RPC method
func (c *Client) GetBlockByNumberRaw(chainID, blockNumber string) (json.RawMessage, error) {
	// Convert blockNumber to proper format if it's "latest"
	tag := blockNumber
	if blockNumber == "latest" {
		tag = "latest"
	} else {
		// If it's a number, convert to hex format (required by eth_getBlockByNumber)
		if _, err := strconv.ParseUint(blockNumber, 10, 64); err == nil {
			blockNum, _ := strconv.ParseUint(blockNumber, 10, 64)
			tag = fmt.Sprintf("0x%x", blockNum)
		}
	}

	params := map[string]string{
		"tag":     tag,
		"boolean": "true", // Include full transaction objects
	}

	return c.Request(chainID, "proxy", "eth_getBlockByNumber", params)
}

// GetBlockRewards gets block rewards by block number
func (c *Client) GetBlockRewards(chainID, blockNumber string) (json.RawMessage, error) {
	params := map[string]string{
		"blockno": blockNumber,
	}

	return c.Request(chainID, "block", "getblockreward", params)
}

// GetContractABI gets the ABI for a verified contract
func (c *Client) GetContractABI(chainID, contractAddress string) (string, error) {
	params := map[string]string{
		"address": contractAddress,
	}

	result, err := c.Request(chainID, "contract", "getabi", params)
	if err != nil {
		return "", err
	}

	var abi string
	if err := json.Unmarshal(result, &abi); err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	return abi, nil
}

// GetContractSourceCode gets the source code of a verified contract
func (c *Client) GetContractSourceCode(chainID, contractAddress string) (json.RawMessage, error) {
	params := map[string]string{
		"address": contractAddress,
	}

	return c.Request(chainID, "contract", "getsourcecode", params)
}

// ExecuteContractMethod executes a read contract function
func (c *Client) ExecuteContractMethod(chainID, contractAddress, methodABI, methodParams string) (json.RawMessage, error) {
	params := map[string]string{
		"to":   contractAddress,
		"data": methodABI,
	}

	if methodParams != "" {
		params["params"] = methodParams
	}

	return c.Request(chainID, "proxy", "eth_call", params)
}

// GetGasOracle gets current gas price oracle output
func (c *Client) GetGasOracle(chainID string) (json.RawMessage, error) {
	return c.Request(chainID, "gastracker", "gasoracle", nil)
}

// GetTokenBalance gets the token balance of an account
func (c *Client) GetTokenBalance(chainID, contractAddress, address string) (string, error) {
	params := map[string]string{
		"contractaddress": contractAddress,
		"address":         address,
		"tag":             "latest",
	}

	result, err := c.Request(chainID, "account", "tokenbalance", params)
	if err != nil {
		return "", err
	}

	var balance string
	if err := json.Unmarshal(result, &balance); err != nil {
		return "", fmt.Errorf("failed to parse token balance: %w", err)
	}

	return balance, nil
}

// GetTransactionByHash gets transaction details by hash
func (c *Client) GetTransactionByHash(chainID, txHash string) (json.RawMessage, error) {
	params := map[string]string{
		"txhash": txHash,
	}

	return c.Request(chainID, "proxy", "eth_getTransactionByHash", params)
}

// GetTransactionByBlockNumberAndIndex gets a transaction by block number and index
func (c *Client) GetTransactionByBlockNumberAndIndex(chainID, blockNumber, index string) (json.RawMessage, error) {
	// Convert blockNumber to proper format if it's not "latest"
	tag := blockNumber
	if blockNumber != "latest" {
		// If it's a number, convert to hex format (required by eth_getTransactionByBlockNumberAndIndex)
		if _, err := strconv.ParseUint(blockNumber, 10, 64); err == nil {
			blockNum, _ := strconv.ParseUint(blockNumber, 10, 64)
			tag = fmt.Sprintf("0x%x", blockNum)
		}
	}

	// Convert index to hex format if it's a number
	idx := index
	if _, err := strconv.ParseUint(index, 10, 64); err == nil {
		indexNum, _ := strconv.ParseUint(index, 10, 64)
		idx = fmt.Sprintf("0x%x", indexNum)
	}

	params := map[string]string{
		"tag":   tag,
		"index": idx,
	}

	return c.Request(chainID, "proxy", "eth_getTransactionByBlockNumberAndIndex", params)
}

// GetTransactionCount gets the number of transactions sent from an address
func (c *Client) GetTransactionCount(chainID, address, tag string) (json.RawMessage, error) {
	if tag == "" {
		tag = "latest"
	}

	params := map[string]string{
		"address": address,
		"tag":     tag,
	}

	return c.Request(chainID, "proxy", "eth_getTransactionCount", params)
}

// GetTransactionReceipt gets transaction receipt
func (c *Client) GetTransactionReceipt(chainID, txHash string) (json.RawMessage, error) {
	params := map[string]string{
		"txhash": txHash,
	}

	return c.Request(chainID, "proxy", "eth_getTransactionReceipt", params)
}

// GetTransactionStatus gets contract execution status for a transaction
func (c *Client) GetTransactionStatus(chainID, txHash string) (json.RawMessage, error) {
	params := map[string]string{
		"txhash": txHash,
	}

	return c.Request(chainID, "transaction", "getstatus", params)
}

// GetTransactionsByAddress gets list of transactions by address
func (c *Client) GetTransactionsByAddress(chainID, address string, params map[string]string) (json.RawMessage, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["address"] = address

	return c.Request(chainID, "account", "txlist", params)
}

// GetInternalTransactionsByAddress gets list of internal transactions by address
func (c *Client) GetInternalTransactionsByAddress(chainID, address string, params map[string]string) (json.RawMessage, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["address"] = address

	return c.Request(chainID, "account", "txlistinternal", params)
}

// GetTokenTransfersByAddress gets list of token transfers by address
func (c *Client) GetTokenTransfersByAddress(chainID, address string, params map[string]string) (json.RawMessage, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["address"] = address

	return c.Request(chainID, "account", "tokentx", params)
}

// GetERC721Transfers gets list of ERC721 token transfers by address
func (c *Client) GetERC721Transfers(chainID, address string, params map[string]string) (json.RawMessage, error) {
	if params == nil {
		params = make(map[string]string)
	}
	params["address"] = address

	return c.Request(chainID, "account", "tokennfttx", params)
}

// TokenDetails represents ERC20 token details
type TokenDetails struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// GetTokenDetails gets comprehensive token information
func (c *Client) GetTokenDetails(chainID, contractAddress string) (json.RawMessage, error) {
	// Handle special addresses for native tokens
	if strings.EqualFold(contractAddress, "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee") {
		// Native ETH/chain token
		nativeSymbol := "ETH" // Default

		// Try to get proper chain symbol based on chainID
		switch chainID {
		case "1": // Ethereum Mainnet
			nativeSymbol = "ETH"
		case "56": // Binance Smart Chain
			nativeSymbol = "BNB"
		case "137": // Polygon
			nativeSymbol = "MATIC"
		case "42161": // Arbitrum One
			nativeSymbol = "ETH"
		case "10": // Optimism
			nativeSymbol = "ETH"
		case "43114": // Avalanche C-Chain
			nativeSymbol = "AVAX"
		case "8453": // Base
			nativeSymbol = "ETH"
		case "324": // zkSync Era
			nativeSymbol = "ETH"
		case "100": // Gnosis
			nativeSymbol = "xDAI"
		case "250": // Fantom
			nativeSymbol = "FTM"
		case "5000": // Mantle
			nativeSymbol = "MNT"
		case "25": // Cronos
			nativeSymbol = "CRO"
		case "1101": // Polygon ZkEVM
			nativeSymbol = "ETH"
		case "59144": // Linea
			nativeSymbol = "ETH"
		case "1284": // Moonbeam
			nativeSymbol = "GLMR"
		case "42220": // Celo
			nativeSymbol = "CELO"
		case "534352": // Scroll
			nativeSymbol = "ETH"
		case "204": // OpBNB
			nativeSymbol = "BNB"
		case "1285": // Moonriver
			nativeSymbol = "MOVR"
		case "42170": // Arbitrum Nova
			nativeSymbol = "ETH"
		case "81457": // Blast
			nativeSymbol = "ETH"
		case "252": // Fraxtal
			nativeSymbol = "frxETH"
		case "1111": // Wemix
			nativeSymbol = "WEMIX"
		case "660279": // Xai
			nativeSymbol = "XAI"
		case "480": // World Chain
			nativeSymbol = "ETH"
		case "33139": // Ape
			nativeSymbol = "APE"
		case "255": // Kroma
			nativeSymbol = "ETH"
		case "167000": // Taiko
			nativeSymbol = "ETH"
		case "199": // Bittorrent
			nativeSymbol = "BTT"
		case "50": // Xdc
			nativeSymbol = "XDC"
		}

		details := TokenDetails{
			Name:     nativeSymbol,
			Symbol:   nativeSymbol,
			Decimals: 18,
		}

		detailsJSON, _ := json.Marshal(details)
		return detailsJSON, nil
	}

	// Special case for USDT on Ethereum
	if chainID == "1" && strings.EqualFold(contractAddress, "0xdAC17F958D2ee523a2206206994597C13D831ec7") {
		details := TokenDetails{
			Name:     "Tether USD",
			Symbol:   "USDT",
			Decimals: 6,
		}

		detailsJSON, _ := json.Marshal(details)
		return detailsJSON, nil
	}

	// Try primary method first - token info endpoint
	params := map[string]string{
		"contractaddress": contractAddress,
	}

	result, err := c.Request(chainID, "token", "tokeninfo", params)
	if err == nil {
		// Check if we got valid token info
		var response map[string]interface{}
		if err := json.Unmarshal(result, &response); err == nil {
			if status, ok := response["status"].(string); ok && status == "1" {
				return result, nil
			}
		}
	}

	// If primary method fails, try with multiple direct contract calls
	// We'll build the details piece by piece
	details := TokenDetails{
		Name:     "Unknown Token",
		Symbol:   "UNKNOWN",
		Decimals: 18, // Default for most ERC20 tokens
	}

	// Try to get token name
	nameResult, _ := c.ExecuteContractMethod(chainID, contractAddress, "0x06fdde03", "") // name()
	if nameResult != nil {
		var hexValue string
		if err := json.Unmarshal(nameResult, &hexValue); err == nil && hexValue != "" {
			// Decode the ABI-encoded string
			details.Name = decodeAbiString(hexValue)
		}
	}

	// Try to get token symbol
	symbolResult, _ := c.ExecuteContractMethod(chainID, contractAddress, "0x95d89b41", "") // symbol()
	if symbolResult != nil {
		var hexValue string
		if err := json.Unmarshal(symbolResult, &hexValue); err == nil && hexValue != "" {
			// Decode the ABI-encoded string
			details.Symbol = decodeAbiString(hexValue)
		}
	}

	// Try to get token decimals
	decimalsResult, _ := c.ExecuteContractMethod(chainID, contractAddress, "0x313ce567", "") // decimals()
	if decimalsResult != nil {
		var hexValue string
		if err := json.Unmarshal(decimalsResult, &hexValue); err == nil && len(hexValue) > 2 {
			if decimals, err := strconv.ParseUint(hexValue[2:], 16, 64); err == nil {
				details.Decimals = int(decimals)
			}
		}
	}

	// Format and return the results
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("error serializing token details: %w", err)
	}

	// Wrap in the standard response format
	responseJSON := fmt.Sprintf(`{"status":"1","message":"OK","result":%s}`, string(detailsJSON))
	return []byte(responseJSON), nil
}

// decodeAbiString decodes an ABI-encoded string from a hex representation
func decodeAbiString(hexData string) string {
	// Format is: 0x + 32 bytes offset + 32 bytes length + data
	if len(hexData) < 2+64+64 {
		// Try treating as bytes32
		if len(hexData) >= 2+64 {
			// Strip 0x prefix
			bytes32Hex := hexData[2:66]
			bytes, err := hexToBytes(bytes32Hex)
			if err == nil {
				// Trim trailing zeros
				return strings.TrimRight(string(bytes), "\x00")
			}
		}
		return ""
	}

	// Skip 0x prefix
	hexData = hexData[2:]

	// Get length (second 32 bytes)
	lengthHex := hexData[64:128]
	length, err := strconv.ParseUint(lengthHex, 16, 64)
	if err != nil {
		return ""
	}

	// If length is 0, return empty string
	if length == 0 {
		return ""
	}

	// Get data (after offset and length)
	if 128+length*2 > uint64(len(hexData)) {
		return ""
	}

	dataHex := hexData[128 : 128+length*2]
	bytes, err := hexToBytes(dataHex)
	if err != nil {
		return ""
	}

	return string(bytes)
}

// hexToBytes converts a hex string to bytes
func hexToBytes(hexStr string) ([]byte, error) {
	length := len(hexStr)
	bytes := make([]byte, length/2)

	for i := 0; i < length; i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		bytes[i/2] = byte(b)
	}

	return bytes, nil
}

// GetLatestBlockNumber gets the latest block number directly
func (c *Client) GetLatestBlockNumber(chainID string) (string, error) {
	result, err := c.Request(chainID, "proxy", "eth_blockNumber", nil)
	if err != nil {
		return "", err
	}

	var blockNumber string
	if err := json.Unmarshal(result, &blockNumber); err != nil {
		return "", fmt.Errorf("failed to parse block number: %w", err)
	}

	// Convert from hex to decimal
	if len(blockNumber) > 2 && blockNumber[:2] == "0x" {
		blockNum, err := strconv.ParseUint(blockNumber[2:], 16, 64)
		if err != nil {
			return blockNumber, nil
		}
		return fmt.Sprintf("%d", blockNum), nil
	}

	return blockNumber, nil
}
