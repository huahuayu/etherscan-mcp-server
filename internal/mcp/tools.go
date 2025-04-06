package mcp

import (
	"context"

	"github.com/huahuayu/etherscan-mcp-server/internal/etherscan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all the Etherscan API tools with the MCP server
func RegisterTools(s *server.MCPServer, client *etherscan.Client) {
	// 1. Get Account Balance
	accountBalanceTool := mcp.NewTool("getAccountBalance",
		mcp.WithDescription("Get the balance of an account on a specific blockchain"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
	)
	s.AddTool(accountBalanceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetAccountBalance(ctx, request, client)
	})

	// 2. Get Block By Number
	blockByNumberTool := mcp.NewTool("getBlockByNumber",
		mcp.WithDescription("Get block information by block number"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("blockNumber",
			mcp.Required(),
			mcp.Description("The block number (or 'latest')"),
		),
	)
	s.AddTool(blockByNumberTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetBlockByNumber(ctx, request, client)
	})

	// 3. Get Block Rewards
	blockRewardsTool := mcp.NewTool("getBlockRewards",
		mcp.WithDescription("Get block rewards by block number"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("blockNumber",
			mcp.Required(),
			mcp.Description("The block number"),
		),
	)
	s.AddTool(blockRewardsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetBlockRewards(ctx, request, client)
	})

	// 4. Get Contract ABI
	contractABITool := mcp.NewTool("getContractABI",
		mcp.WithDescription("Get the ABI for a verified contract"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("contractAddress",
			mcp.Required(),
			mcp.Description("The contract address"),
		),
	)
	s.AddTool(contractABITool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetContractABI(ctx, request, client)
	})

	// 5. Get Contract Source Code
	contractSourceCodeTool := mcp.NewTool("getContractSourceCode",
		mcp.WithDescription("Get the source code of a verified contract"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("contractAddress",
			mcp.Required(),
			mcp.Description("The contract address"),
		),
	)
	s.AddTool(contractSourceCodeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetContractSourceCode(ctx, request, client)
	})

	// 6. Execute Contract Method
	executeContractMethodTool := mcp.NewTool("executeContractMethod",
		mcp.WithDescription("Execute a read contract function"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("contractAddress",
			mcp.Required(),
			mcp.Description("The contract address"),
		),
		mcp.WithString("methodABI",
			mcp.Required(),
			mcp.Description("The ABI of the contract method"),
		),
		mcp.WithString("methodParams",
			mcp.Description("Comma-separated parameter values for the method"),
		),
	)
	s.AddTool(executeContractMethodTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleExecuteContractMethod(ctx, request, client)
	})

	// 7. Get Gas Oracle
	gasOracleTool := mcp.NewTool("getGasOracle",
		mcp.WithDescription("Get current gas price oracle output"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
	)
	s.AddTool(gasOracleTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetGasOracle(ctx, request, client)
	})

	// 8. Get Token Balance
	tokenBalanceTool := mcp.NewTool("getTokenBalance",
		mcp.WithDescription("Get the token balance of an account on a specific blockchain"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("contractAddress",
			mcp.Required(),
			mcp.Description("The token contract address"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
	)
	s.AddTool(tokenBalanceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTokenBalance(ctx, request, client)
	})

	// 9. Get Token Details
	tokenDetailsTool := mcp.NewTool("getTokenDetails",
		mcp.WithDescription("Get comprehensive token information"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("contractAddress",
			mcp.Required(),
			mcp.Description("The token contract address"),
		),
	)
	s.AddTool(tokenDetailsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTokenDetails(ctx, request, client)
	})

	// 10. Get Transaction By Hash
	transactionByHashTool := mcp.NewTool("getTransactionByHash",
		mcp.WithDescription("Get transaction details by hash"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("txHash",
			mcp.Required(),
			mcp.Description("The transaction hash"),
		),
	)
	s.AddTool(transactionByHashTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionByHash(ctx, request, client)
	})

	// 10a. Get Transaction By Block Number And Index
	transactionByBlockNumberAndIndexTool := mcp.NewTool("getTransactionByBlockNumberAndIndex",
		mcp.WithDescription("Get transaction by block number and index"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("blockNumber",
			mcp.Required(),
			mcp.Description("The block number (or 'latest')"),
		),
		mcp.WithString("index",
			mcp.Required(),
			mcp.Description("The transaction index position"),
		),
	)
	s.AddTool(transactionByBlockNumberAndIndexTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionByBlockNumberAndIndex(ctx, request, client)
	})

	// 10b. Get Transaction Count
	transactionCountTool := mcp.NewTool("getTransactionCount",
		mcp.WithDescription("Get the number of transactions sent from an address"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
		mcp.WithString("tag",
			mcp.Description("The block number tag (default: 'latest')"),
		),
	)
	s.AddTool(transactionCountTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionCount(ctx, request, client)
	})

	// 10c. Get Transaction Receipt Status
	transactionReceiptTool := mcp.NewTool("getTransactionReceipt",
		mcp.WithDescription("Check transaction receipt status"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("txHash",
			mcp.Required(),
			mcp.Description("The transaction hash"),
		),
	)
	s.AddTool(transactionReceiptTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionReceipt(ctx, request, client)
	})

	// 10c. Get Transaction Status
	transactionStatusTool := mcp.NewTool("getTransactionStatus",
		mcp.WithDescription("Check contract execution status"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("txHash",
			mcp.Required(),
			mcp.Description("The transaction hash"),
		),
	)
	s.AddTool(transactionStatusTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionStatus(ctx, request, client)
	})

	// 11. Get Transactions By Address
	transactionsByAddressTool := mcp.NewTool("getTransactionsByAddress",
		mcp.WithDescription("Get list of transactions by address"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
		mcp.WithString("startBlock",
			mcp.Description("Starting block number"),
		),
		mcp.WithString("endBlock",
			mcp.Description("Ending block number"),
		),
		mcp.WithString("page",
			mcp.Description("Page number"),
		),
		mcp.WithString("offset",
			mcp.Description("Number of records to return"),
		),
	)
	s.AddTool(transactionsByAddressTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTransactionsByAddress(ctx, request, client)
	})

	// 12. Get Internal Transactions By Address
	internalTransactionsByAddressTool := mcp.NewTool("getInternalTransactionsByAddress",
		mcp.WithDescription("Get list of internal transactions by address"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
		mcp.WithString("startBlock",
			mcp.Description("Starting block number"),
		),
		mcp.WithString("endBlock",
			mcp.Description("Ending block number"),
		),
		mcp.WithString("page",
			mcp.Description("Page number"),
		),
		mcp.WithString("offset",
			mcp.Description("Number of records to return"),
		),
	)
	s.AddTool(internalTransactionsByAddressTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetInternalTransactionsByAddress(ctx, request, client)
	})

	// 13. Get Token Transfers By Address
	tokenTransfersByAddressTool := mcp.NewTool("getTokenTransfersByAddress",
		mcp.WithDescription("Get list of token transfers by address"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
		mcp.WithString("contractAddress",
			mcp.Description("The token contract address"),
		),
		mcp.WithString("startBlock",
			mcp.Description("Starting block number"),
		),
		mcp.WithString("endBlock",
			mcp.Description("Ending block number"),
		),
		mcp.WithString("page",
			mcp.Description("Page number"),
		),
		mcp.WithString("offset",
			mcp.Description("Number of records to return"),
		),
	)
	s.AddTool(tokenTransfersByAddressTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetTokenTransfersByAddress(ctx, request, client)
	})

	// 14. Get ERC721 Transfers
	erc721TransfersTool := mcp.NewTool("getERC721Transfers",
		mcp.WithDescription("Get list of ERC721 token transfers by address"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum)"),
		),
		mcp.WithString("address",
			mcp.Required(),
			mcp.Description("The account address"),
		),
		mcp.WithString("contractAddress",
			mcp.Description("The token contract address"),
		),
		mcp.WithString("startBlock",
			mcp.Description("Starting block number"),
		),
		mcp.WithString("endBlock",
			mcp.Description("Ending block number"),
		),
		mcp.WithString("page",
			mcp.Description("Page number"),
		),
		mcp.WithString("offset",
			mcp.Description("Number of records to return"),
		),
	)
	s.AddTool(erc721TransfersTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetERC721Transfers(ctx, request, client)
	})

	// Add new tool registration for getLatestBlockNumber
	latestBlockNumberTool := mcp.NewTool("getLatestBlockNumber",
		mcp.WithDescription("Get the latest block number"),
		mcp.WithString("chainID",
			mcp.Required(),
			mcp.Description("The chain ID (e.g., 1 for Ethereum, 42161 for Arbitrum)"),
		),
	)
	s.AddTool(latestBlockNumberTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleGetLatestBlockNumber(ctx, request, client)
	})
}
