.PHONY: build clean install

# Default target
all: build

# Build the server
build:
	@echo "Building etherscan-mcp-server..."
	@go build -o bin/etherscan-mcp-server ./cmd/etherscan-mcp-server
	@chmod +x bin/etherscan-mcp-server
	@echo "Build complete: bin/etherscan-mcp-server"

# Install the server to /usr/local/bin
install: build
	@echo "Installing etherscan-mcp-server..."
	@mkdir -p /usr/local/bin
	@cp bin/etherscan-mcp-server /usr/local/bin/
	@echo "Installation complete: /usr/local/bin/etherscan-mcp-server"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/etherscan-mcp-server
	@echo "Clean complete" 