package mcp

import (
	"context"
	"fmt"

	"github.com/huahuayu/etherscan-mcp-server/internal/etherscan"
	"github.com/mark3labs/mcp-go/mcp"
)

// Handler functions
func handleGetAccountBalance(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	balance, err := client.GetAccountBalance(chainID, address)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"balance": "%s"}`, balance)), nil
}

func handleGetBlockByNumber(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	blockNumber, ok := request.Params.Arguments["blockNumber"].(string)
	if !ok {
		return nil, fmt.Errorf("blockNumber must be a string")
	}

	result, err := client.GetBlockByNumber(chainID, blockNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetBlockRewards(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	blockNumber, ok := request.Params.Arguments["blockNumber"].(string)
	if !ok {
		return nil, fmt.Errorf("blockNumber must be a string")
	}

	result, err := client.GetBlockRewards(chainID, blockNumber)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetContractABI(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	contractAddress, ok := request.Params.Arguments["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	abi, err := client.GetContractABI(chainID, contractAddress)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(abi), nil
}

func handleGetContractSourceCode(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	contractAddress, ok := request.Params.Arguments["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	result, err := client.GetContractSourceCode(chainID, contractAddress)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleExecuteContractMethod(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	contractAddress, ok := request.Params.Arguments["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	methodABI, ok := request.Params.Arguments["methodABI"].(string)
	if !ok {
		return nil, fmt.Errorf("methodABI must be a string")
	}

	methodParams, _ := request.Params.Arguments["methodParams"].(string)

	result, err := client.ExecuteContractMethod(chainID, contractAddress, methodABI, methodParams)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetGasOracle(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	result, err := client.GetGasOracle(chainID)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTokenBalance(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	contractAddress, ok := request.Params.Arguments["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	balance, err := client.GetTokenBalance(chainID, contractAddress, address)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"balance": "%s"}`, balance)), nil
}

func handleGetTokenDetails(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	contractAddress, ok := request.Params.Arguments["contractAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("contractAddress must be a string")
	}

	result, err := client.GetTokenDetails(chainID, contractAddress)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTransactionByHash(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	txHash, ok := request.Params.Arguments["txHash"].(string)
	if !ok {
		return nil, fmt.Errorf("txHash must be a string")
	}

	result, err := client.GetTransactionByHash(chainID, txHash)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTransactionReceipt(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	txHash, ok := request.Params.Arguments["txHash"].(string)
	if !ok {
		return nil, fmt.Errorf("txHash must be a string")
	}

	result, err := client.GetTransactionReceipt(chainID, txHash)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTransactionStatus(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	txHash, ok := request.Params.Arguments["txHash"].(string)
	if !ok {
		return nil, fmt.Errorf("txHash must be a string")
	}

	result, err := client.GetTransactionStatus(chainID, txHash)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTransactionsByAddress(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	params := make(map[string]string)

	if startBlock, ok := request.Params.Arguments["startBlock"].(string); ok && startBlock != "" {
		params["startblock"] = startBlock
	}

	if endBlock, ok := request.Params.Arguments["endBlock"].(string); ok && endBlock != "" {
		params["endblock"] = endBlock
	}

	if page, ok := request.Params.Arguments["page"].(string); ok && page != "" {
		params["page"] = page
	}

	if offset, ok := request.Params.Arguments["offset"].(string); ok && offset != "" {
		params["offset"] = offset
	}

	result, err := client.GetTransactionsByAddress(chainID, address, params)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetInternalTransactionsByAddress(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	params := make(map[string]string)

	if startBlock, ok := request.Params.Arguments["startBlock"].(string); ok && startBlock != "" {
		params["startblock"] = startBlock
	}

	if endBlock, ok := request.Params.Arguments["endBlock"].(string); ok && endBlock != "" {
		params["endblock"] = endBlock
	}

	if page, ok := request.Params.Arguments["page"].(string); ok && page != "" {
		params["page"] = page
	}

	if offset, ok := request.Params.Arguments["offset"].(string); ok && offset != "" {
		params["offset"] = offset
	}

	result, err := client.GetInternalTransactionsByAddress(chainID, address, params)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTokenTransfersByAddress(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	params := make(map[string]string)

	if contractAddress, ok := request.Params.Arguments["contractAddress"].(string); ok && contractAddress != "" {
		params["contractaddress"] = contractAddress
	}

	if startBlock, ok := request.Params.Arguments["startBlock"].(string); ok && startBlock != "" {
		params["startblock"] = startBlock
	}

	if endBlock, ok := request.Params.Arguments["endBlock"].(string); ok && endBlock != "" {
		params["endblock"] = endBlock
	}

	if page, ok := request.Params.Arguments["page"].(string); ok && page != "" {
		params["page"] = page
	}

	if offset, ok := request.Params.Arguments["offset"].(string); ok && offset != "" {
		params["offset"] = offset
	}

	result, err := client.GetTokenTransfersByAddress(chainID, address, params)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetERC721Transfers(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	params := make(map[string]string)

	if contractAddress, ok := request.Params.Arguments["contractAddress"].(string); ok && contractAddress != "" {
		params["contractaddress"] = contractAddress
	}

	if startBlock, ok := request.Params.Arguments["startBlock"].(string); ok && startBlock != "" {
		params["startblock"] = startBlock
	}

	if endBlock, ok := request.Params.Arguments["endBlock"].(string); ok && endBlock != "" {
		params["endblock"] = endBlock
	}

	if page, ok := request.Params.Arguments["page"].(string); ok && page != "" {
		params["page"] = page
	}

	if offset, ok := request.Params.Arguments["offset"].(string); ok && offset != "" {
		params["offset"] = offset
	}

	result, err := client.GetERC721Transfers(chainID, address, params)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetLatestBlockNumber(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	blockNumber, err := client.GetLatestBlockNumber(chainID)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"blockNumber": "%s"}`, blockNumber)), nil
}

func handleGetTransactionByBlockNumberAndIndex(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	blockNumber, ok := request.Params.Arguments["blockNumber"].(string)
	if !ok {
		return nil, fmt.Errorf("blockNumber must be a string")
	}

	index, ok := request.Params.Arguments["index"].(string)
	if !ok {
		return nil, fmt.Errorf("index must be a string")
	}

	result, err := client.GetTransactionByBlockNumberAndIndex(chainID, blockNumber, index)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetTransactionCount(ctx context.Context, request mcp.CallToolRequest, client *etherscan.Client) (*mcp.CallToolResult, error) {
	chainID, ok := request.Params.Arguments["chainID"].(string)
	if !ok {
		return nil, fmt.Errorf("chainID must be a string")
	}

	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, fmt.Errorf("address must be a string")
	}

	tag, _ := request.Params.Arguments["tag"].(string)

	result, err := client.GetTransactionCount(chainID, address, tag)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(result)), nil
}
