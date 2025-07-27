#!/bin/bash

# Weather MCP Server Startup Script
# This script starts the weather MCP server with the API key

echo "Starting Weather MCP Server..."
echo "API Key: f5fc1e71ee4a4963a74195621252607"
echo "Server supports international locations like Tokyo, London, New York, etc."
echo ""

# Set the API key
export WEATHER_API_KEY=f5fc1e71ee4a4963a74195621252607

# Change to the weather-mcp-server directory
cd weather-mcp-server

# Start the server in stdio mode (for MCP)
echo "Server is ready for MCP connections..."
echo "Use this path in your Claude MCP configuration:"
echo "$(pwd)/weather-mcp-server"
echo ""

# Run the server
./weather-mcp-server 