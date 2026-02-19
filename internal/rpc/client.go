package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Pre-configured free RPC endpoints from LlamaRPC
var chainRPCURLs = map[string]string{
	"56":    "https://binance.llamarpc.com",          // BSC
	"8453":  "https://base.llamarpc.com",             // Base
	"43114": "https://api.avax.network/ext/bc/C/rpc", // Avalanche C-Chain
}

// IsRPCFallbackChain checks if a chain has RPC fallback support
func IsRPCFallbackChain(chainID string) bool {
	_, ok := chainRPCURLs[chainID]
	return ok
}

// Client is a JSON-RPC client for direct RPC calls
type Client struct {
	httpClient *http.Client
}

// jsonRPCRequest represents a JSON-RPC 2.0 request
type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// jsonRPCResponse represents a JSON-RPC 2.0 response
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

// jsonRPCError represents a JSON-RPC error
type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new RPC client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// call performs a JSON-RPC call to the appropriate chain RPC endpoint
func (c *Client) call(chainID, method string, params []interface{}) (json.RawMessage, error) {
	rpcURL, ok := chainRPCURLs[chainID]
	if !ok {
		return nil, fmt.Errorf("no RPC endpoint configured for chain %s", chainID)
	}

	reqBody := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RPC request: %w", err)
	}

	req, err := http.NewRequest("POST", rpcURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RPC request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RPC response: %w", err)
	}

	var rpcResp jsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse RPC response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// BlockNumber returns the latest block number (decimal string)
func (c *Client) BlockNumber(chainID string) (string, error) {
	result, err := c.call(chainID, "eth_blockNumber", []interface{}{})
	if err != nil {
		return "", err
	}

	var hexBlock string
	if err := json.Unmarshal(result, &hexBlock); err != nil {
		return "", fmt.Errorf("failed to parse block number: %w", err)
	}

	// Convert hex to decimal
	if len(hexBlock) > 2 && hexBlock[:2] == "0x" {
		blockNum, err := strconv.ParseUint(hexBlock[2:], 16, 64)
		if err != nil {
			return hexBlock, nil
		}
		return fmt.Sprintf("%d", blockNum), nil
	}

	return hexBlock, nil
}

// GetBalance returns the balance of an address in wei (decimal string)
func (c *Client) GetBalance(chainID, address string) (string, error) {
	result, err := c.call(chainID, "eth_getBalance", []interface{}{address, "latest"})
	if err != nil {
		return "", err
	}

	var hexBalance string
	if err := json.Unmarshal(result, &hexBalance); err != nil {
		return "", fmt.Errorf("failed to parse balance: %w", err)
	}

	// Convert hex to decimal
	if len(hexBalance) > 2 && hexBalance[:2] == "0x" {
		balance := new(big.Int)
		balance.SetString(hexBalance[2:], 16)
		return balance.String(), nil
	}

	return hexBalance, nil
}

// GetTokenBalance returns the ERC20 token balance of an address (decimal string)
func (c *Client) GetTokenBalance(chainID, contractAddress, address string) (string, error) {
	// balanceOf(address) selector = 0x70a08231
	// Pad address to 32 bytes
	paddedAddress := fmt.Sprintf("0x70a08231%064s", strings.TrimPrefix(address, "0x"))
	paddedAddress = strings.ReplaceAll(paddedAddress, " ", "0")

	callData := map[string]string{
		"to":   contractAddress,
		"data": paddedAddress,
	}

	result, err := c.call(chainID, "eth_call", []interface{}{callData, "latest"})
	if err != nil {
		return "", err
	}

	var hexBalance string
	if err := json.Unmarshal(result, &hexBalance); err != nil {
		return "", fmt.Errorf("failed to parse token balance: %w", err)
	}

	// Convert hex to decimal
	if len(hexBalance) > 2 && hexBalance[:2] == "0x" {
		balance := new(big.Int)
		balance.SetString(hexBalance[2:], 16)
		return balance.String(), nil
	}

	return hexBalance, nil
}

// TokenDetails represents ERC20 token details
type TokenDetails struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// GetTokenDetails returns ERC20 token name, symbol, and decimals
func (c *Client) GetTokenDetails(chainID, contractAddress string) (json.RawMessage, error) {
	details := TokenDetails{
		Name:     "Unknown Token",
		Symbol:   "UNKNOWN",
		Decimals: 18,
	}

	// name() = 0x06fdde03
	nameResult, _ := c.EthCall(chainID, contractAddress, "0x06fdde03")
	if nameResult != nil {
		var hexValue string
		if err := json.Unmarshal(nameResult, &hexValue); err == nil && hexValue != "" {
			if decoded := decodeAbiString(hexValue); decoded != "" {
				details.Name = decoded
			}
		}
	}

	// symbol() = 0x95d89b41
	symbolResult, _ := c.EthCall(chainID, contractAddress, "0x95d89b41")
	if symbolResult != nil {
		var hexValue string
		if err := json.Unmarshal(symbolResult, &hexValue); err == nil && hexValue != "" {
			if decoded := decodeAbiString(hexValue); decoded != "" {
				details.Symbol = decoded
			}
		}
	}

	// decimals() = 0x313ce567
	decimalsResult, _ := c.EthCall(chainID, contractAddress, "0x313ce567")
	if decimalsResult != nil {
		var hexValue string
		if err := json.Unmarshal(decimalsResult, &hexValue); err == nil && len(hexValue) > 2 {
			if decimals, err := strconv.ParseUint(hexValue[2:], 16, 64); err == nil {
				details.Decimals = int(decimals)
			}
		}
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("error serializing token details: %w", err)
	}

	responseJSON := fmt.Sprintf(`{"status":"1","message":"OK","result":%s}`, string(detailsJSON))
	return []byte(responseJSON), nil
}

// GetTransactionByHash returns transaction details by hash
func (c *Client) GetTransactionByHash(chainID, txHash string) (json.RawMessage, error) {
	return c.call(chainID, "eth_getTransactionByHash", []interface{}{txHash})
}

// GetTransactionReceipt returns the transaction receipt
func (c *Client) GetTransactionReceipt(chainID, txHash string) (json.RawMessage, error) {
	return c.call(chainID, "eth_getTransactionReceipt", []interface{}{txHash})
}

// GetTransactionCount returns the number of transactions from an address
func (c *Client) GetTransactionCount(chainID, address, tag string) (json.RawMessage, error) {
	if tag == "" {
		tag = "latest"
	}
	return c.call(chainID, "eth_getTransactionCount", []interface{}{address, tag})
}

// EthCall performs a read-only contract call
func (c *Client) EthCall(chainID, to, data string) (json.RawMessage, error) {
	callData := map[string]string{
		"to":   to,
		"data": data,
	}
	return c.call(chainID, "eth_call", []interface{}{callData, "latest"})
}

// decodeAbiString decodes an ABI-encoded string from a hex representation
func decodeAbiString(hexData string) string {
	// Format is: 0x + 32 bytes offset + 32 bytes length + data
	if len(hexData) < 2+64+64 {
		// Try treating as bytes32
		if len(hexData) >= 2+64 {
			bytes32Hex := hexData[2:66]
			bytes, err := hexToBytes(bytes32Hex)
			if err == nil {
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

	if length == 0 {
		return ""
	}

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
	result := make([]byte, length/2)

	for i := 0; i < length; i += 2 {
		b, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		result[i/2] = byte(b)
	}

	return result, nil
}
