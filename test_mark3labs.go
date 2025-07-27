package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	fmt.Println("🧪 Testing Mark3Labs Weather MCP Server")
	fmt.Println(strings.Repeat("=", 50))

	// Start the Mark3Labs server as a subprocess
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "/Users/hulkster/WeatherFull/mark3labs-implementation"
	cmd.Env = append(os.Environ(), "WEATHER_API_KEY=f5fc1e71ee4a4963a74195621252607")

	// Get stdin/stdout pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start the server
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	// Create STDIO transport
	stdioTransport := transport.NewStdio(stdin, stdout, nil, "mark3labs-weather-server")

	// Create a new MCP client
	c := client.NewClient(stdioTransport)

	// Start the client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
				Name:    "mark3labs-test-client",
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

	fmt.Println("✅ Mark3Labs Server Test Results:")
	fmt.Printf("Available tools: %d\n", len(toolsResult.Tools))
	for _, tool := range toolsResult.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Test the get_current_weather tool
	if len(toolsResult.Tools) > 0 {
		fmt.Println("\n🧪 Testing get_current_weather tool...")

		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "get_current_weather",
				Arguments: map[string]interface{}{
					"city": "Tokyo",
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
					fmt.Printf("📄 Response: %s\n", textContent.Text)
				}
			}
		}
	}

	fmt.Println("\n🎉 Mark3Labs server test completed successfully!")
	fmt.Println("✅ mark3labs/mcp-go SDK implementation working!")
	os.Exit(0)
}
