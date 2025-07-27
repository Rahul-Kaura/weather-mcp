package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	fmt.Println("🧪 Testing TuanKiri Weather MCP Server")
	fmt.Println(strings.Repeat("=", 50))

	// Create SSE transport
	sseTransport, err := transport.NewSSE("http://localhost:8000/sse")
	if err != nil {
		log.Fatalf("Failed to create SSE transport: %v", err)
	}

	// Create a new MCP client
	c := client.NewClient(sseTransport)

	// Start the client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer c.Close()

	// Initialize the client
	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities: mcp.ClientCapabilities{
				Experimental: map[string]any{
					"tools": true,
				},
			},
			ClientInfo: mcp.Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
	}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// List available tools
	toolsRequest := mcp.ListToolsRequest{}
	toolsResult, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Println("✅ TuanKiri Server Test Results:")
	fmt.Printf("Available tools: %d\n", len(toolsResult.Tools))
	for _, tool := range toolsResult.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Test the current_weather tool
	if len(toolsResult.Tools) > 0 {
		fmt.Println("\n🧪 Testing current_weather tool...")
		
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "current_weather",
				Arguments: map[string]interface{}{
					"city": "London",
				},
			},
		}

		result, err := c.CallTool(ctx, request)
		if err != nil {
			log.Printf("❌ Failed to call tool: %v", err)
		} else {
			fmt.Println("✅ Tool call successful!")
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(mcp.TextContent); ok {
					fmt.Printf("📄 Response length: %d characters\n", len(textContent.Text))
				}
			}
		}
	}

	fmt.Println("\n🎉 TuanKiri server test completed successfully!")
	os.Exit(0)
}
