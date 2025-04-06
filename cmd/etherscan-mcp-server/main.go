package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huahuayu/etherscan-mcp-server/internal/etherscan"
	"github.com/huahuayu/etherscan-mcp-server/internal/mcp"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	// Get environment variables with defaults
	apiKey := getEnv("ETHERSCAN_API_KEY", "")
	if apiKey == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable is required")
	}

	port := getEnv("PORT", "4000")
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

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"Etherscan MCP Server",
		"1.0.0",
	)

	// Register tools
	mcp.RegisterTools(mcpServer, client)

	// Create custom SSE server
	sseServer := mcp.NewCustomSSEServer(mcpServer)
	sseServer.WithHeartbeatInterval(25 * time.Second)

	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting Etherscan MCP Server on %s", addr)
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
