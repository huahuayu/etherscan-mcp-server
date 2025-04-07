package mcp

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
)

// CustomStdioServer extends the StdioServer with any custom functionality
type CustomStdioServer struct {
	server *server.StdioServer
}

// NewCustomStdioServer creates a new CustomStdioServer
func NewCustomStdioServer(mcpServer *server.MCPServer) *CustomStdioServer {
	return &CustomStdioServer{
		server: server.NewStdioServer(mcpServer),
	}
}

// Start begins serving on stdin/stdout
func (s *CustomStdioServer) Start() error {
	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start listening on stdin/stdout
	return s.server.Listen(ctx, os.Stdin, os.Stdout)
}

// SetErrorLogger configures where error messages are logged
func (s *CustomStdioServer) SetErrorLogger(logger *log.Logger) {
	s.server.SetErrorLogger(logger)
}
