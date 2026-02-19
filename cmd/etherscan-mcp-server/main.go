package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huahuayu/etherscan-mcp-server/internal/etherscan"
	"github.com/huahuayu/etherscan-mcp-server/internal/mcp"
	"github.com/huahuayu/etherscan-mcp-server/internal/rpc"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Define flags
	useSSE := flag.Bool("sse", false, "Use SSE server mode (default is stdin/stdout)")
	port := flag.String("port", "", "Port for SSE server (defaults to PORT env var or 4000)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	// Get environment variables with defaults
	apiKey := getEnv("ETHERSCAN_API_KEY", "")
	if apiKey == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable is required")
	}

	// Check env var for SSE mode (overrides flag if set)
	useSSEEnv := getEnv("USE_SSE", "")
	if useSSEEnv == "true" {
		*useSSE = true
	} else if useSSEEnv == "false" {
		*useSSE = false
	}

	// If port flag not set, get from env or use default
	if *port == "" {
		*port = getEnv("PORT", "4000")
	}

	logLevel := getEnv("LOG_LEVEL", "info")

	// Configure logging
	switch logLevel {
	case "debug":
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	default:
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// Initialize Etherscan client
	client := etherscan.NewClient(apiKey)

	// Initialize RPC client for fallback
	rpcClient := rpc.NewClient()

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"Etherscan MCP Server",
		"1.0.0",
	)

	// Register tools
	mcp.RegisterTools(mcpServer, client, rpcClient)

	if *useSSE {
		// SSE server mode
		log.Printf("Starting in SSE mode...")
		runSSEServer(mcpServer, *port)
	} else {
		// Default StdIO server mode
		log.Printf("Starting in StdIO mode...")
		runStdIOServer(mcpServer)
	}
}

func runSSEServer(mcpServer *server.MCPServer, port string) {
	// Create custom SSE server
	sseServer := mcp.NewCustomSSEServer(mcpServer)
	sseServer.WithHeartbeatInterval(25 * time.Second)

	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting Etherscan MCP Server in SSE mode on %s", addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interruption signal
	<-ctx.Done()
	stop()
	log.Println("Shutting down server...")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := sseServer.GracefulShutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func runStdIOServer(mcpServer *server.MCPServer) {
	// Create custom StdIO server
	stdioServer := mcp.NewCustomStdioServer(mcpServer)

	// Configure error logger to write to stderr
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	stdioServer.SetErrorLogger(errorLogger)

	log.Printf("Starting Etherscan MCP Server in StdIO mode")

	// Start listening on stdin/stdout
	if err := stdioServer.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
