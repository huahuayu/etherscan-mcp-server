# Etherscan MCP Server

A Go implementation of an Etherscan API client for the [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go) that enables LLM applications to access Etherscan blockchain data.

## Features

- Access Etherscan API V2 for multi-chain support
- Use a single API key for over [50 supported chains](https://docs.etherscan.io/etherscan-v2/getting-started/supported-chains)
- Supports various blockchain data retrieval methods:
  - Account balances
  - Block information
  - Contract data (ABI, source code)
  - Gas oracle
  - Token information
  - Transaction data
  - And more

## Requirements

- Etherscan API key (get from https://etherscan.io/myapikey)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/huahuayu/etherscan-mcp-server.git
   cd etherscan-mcp-server
   ```

2. Set Etherscan API key env variable:
   ```bash
   ETHERSCAN_API_KEY=your_api_key_here
   ```

3. Build the server:
   ```bash
   make build
   ```

4. Install the server:
   ```bash
   make install # install to /usr/local/bin
   ```

## Usage

### Default Mode (Standard Input/Output)

Run the server using stdin/stdout communication (default mode):
```bash
./bin/etherscan-mcp-server
```

This mode is useful for direct integration with LLM applications that communicate via stdin/stdout.

MCP config:
```json
{
  "mcpServers": {
    "etherscan-mcp-server": {
      "command": "etherscan-mcp-server",
      "env": {
        "ETHERSCAN_API_KEY": "$your_api_key"
      }
    }
   }
}
```

Restart cursor and check if it's success:

![](https://cdn.0xbuilder.com/img/20250407121209988.png)


### SSE Mode

Run the server in Server-Sent Events mode:
```bash
./bin/etherscan-mcp-server --sse
```

In SSE mode, the server listens on HTTP and provides an SSE endpoint.

SSE MCP config:
```json
{
  "mcpServers": {
    "etherscan-mcp-server": {
      "url": "http://localhost:4000/sse",
      "env": {
        "ETHERSCAN_API_KEY": "$your_api_key"
      }
    }
   }
}
```

#### Server Options

- `--sse`: Enable SSE server mode (default is stdin/stdout mode)
- `--port <port>`: Specify the port for SSE server (defaults to PORT env var or 4000)

### Connection Endpoints (SSE Mode)

When running in SSE mode, the server provides:

- **SSE Endpoint**: `http://localhost:4000/sse`

## Supported Chains

50+ chains supported, please refer to the [Etherscan V2 Supported Chains](https://docs.etherscan.io/supported-chains) for a complete list and their corresponding Chain IDs.

> **Note on Free Tier Limitations:**
> According to Etherscan's [supported chains policy](https://docs.etherscan.io/supported-chains), some networks are **not available** under the Free API plan. These include (but are not limited to):
> - **BNB Smart Chain (BSC)**
> - **Base**
> - **Avalanche C-Chain**
>
> To access these chains via the Etherscan V2 API, a paid tier plan is required.

## Example Queries

You can use natural language queries like these:

### Account and Balance Information
- "What's the ETH balance of address 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae?"
- "Show me the token balance for USDT on address 0x123abc... on BSC"
- "How many transactions has 0xvitalik.eth made from this address?"

### Block Information
- "Get information about the latest Polygon block"
- "What are the rewards for miners in block 17000000?"
- "Who mined block 16900000 on Ethereum?"

### Contract Interaction
- "Show me the source code at 0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"
- "What's the ABI for the USDC contract on Polygon?"
- "Call the balanceOf function of USDT contract with my address as parameter"

### Transaction Information
- "Give me details for transaction 0x123456789abcdef..."
- "Has transaction 0xabcdef... been confirmed yet?"
- "What was the gas price used in transaction 0x789abc..."

### Gas and Network
- "What are the current gas prices on Ethereum?"
- "What's the recommended gas price for a fast transaction right now?"
- "How many transactions are pending on Ethereum network?"

### Token Information
- "Tell me about the LINK token contract"
- "What ERC-721 NFTs does address 0x123... own?"
- "Show recent token transfers for 0xvitalik.eth"

### Custom Queries
- "Track all USDC transfers to Binance hot wallet in the last 1000 blocks"
- "Which addresses received the most ETH in block 17000000?"
- "Compare gas usage on Ethereum vs Arbitrum for similar transactions"

## Available Tools

The Etherscan MCP Server provides the following tools for accessing blockchain data:

1. **getAccountBalance** - Get the balance of an account on a specific blockchain
2. **getBlockByNumber** - Get block information by block number
3. **getBlockRewards** - Get block rewards by block number
4. **getContractABI** - Get the ABI for a verified contract
5. **getContractSourceCode** - Get the source code of a verified contract
6. **executeContractMethod** - Execute a read contract function
7. **getGasOracle** - Get current gas price oracle output
8. **getTokenBalance** - Get the token balance of an account on a specific blockchain
9. **getTokenDetails** - Get comprehensive token information
10. **getTransactionByHash** - Get transaction details by hash
11. **getTransactionByBlockNumberAndIndex** - Get transaction by block number and index
12. **getTransactionCount** - Get the number of transactions sent from an address
13. **getTransactionReceipt** - Check transaction receipt status
14. **getTransactionStatus** - Check contract execution status
15. **getTransactionsByAddress** - Get list of transactions by address
16. **getInternalTransactionsByAddress** - Get list of internal transactions by address
17. **getTokenTransfersByAddress** - Get list of token transfers by address
18. **getERC721Transfers** - Get list of ERC721 token transfers by address
19. **getLatestBlockNumber** - Get the latest block number

Each tool accepts specific parameters and provides blockchain data in a structured format.

## License

MIT