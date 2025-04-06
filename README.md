# Etherscan MCP Server

A Go implementation of an Etherscan API client for the [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go) that enables LLM applications to access Etherscan blockchain data.

## Features

- Access Etherscan API V2 for multi-chain support
- Use a single API key for over 50 supported chains
- Supports various blockchain data retrieval methods:
  - Account balances
  - Block information
  - Contract data (ABI, source code)
  - Gas oracle
  - Token information
  - Transaction data
  - And more

## Requirements

- Etherscan API key

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
   chmod +x build.sh
   ./build.sh
   ```

## Usage

Run the server:
```bash
./bin/etherscan-mcp-server
```

The server will start and listen on the port specified in your .env file (default: 4000).

### Connection Methods

The server supports Server-Sent Events (SSE) connections:

- **SSE Endpoint**: `http://localhost:4000/sse`

## Supported Chains

50+ chains supported, plz refer https://docs.etherscan.io/etherscan-v2/getting-started/supported-chains

## Example Queries

You can use natural language queries like these:

### Account and Balance Information
- "What's the ETH balance of address 0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae?"
- "Show me the token balance for USDT on address 0x123abc... on Ethereum"
- "How many transactions has 0xvitalik.eth made from this address?"

### Block Information
- "Get information about the latest Ethereum block"
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

## License

MIT