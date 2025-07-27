#!/usr/bin/env python3
"""
Simple test script for the weather MCP server
"""
import subprocess
import json
import sys

def test_weather_server():
    """Test the weather server with international cities"""
    
    # Test cities
    test_cities = ["Tokyo", "London", "New York", "Paris", "Sydney"]
    
    print("Testing Weather MCP Server with international cities...")
    print("=" * 50)
    
    for city in test_cities:
        print(f"\nTesting: {city}")
        try:
            # This would be the actual MCP call in a real implementation
            print(f"  ✓ Would get weather for {city}")
        except Exception as e:
            print(f"  ✗ Error getting weather for {city}: {e}")
    
    print("\n" + "=" * 50)
    print("To actually test this, you need to:")
    print("1. Get a WeatherAPI.com API key")
    print("2. Set WEATHER_API_KEY environment variable")
    print("3. Run: ./weather-mcp-server/weather-mcp-server")
    print("4. Use an MCP client to connect to the server")

if __name__ == "__main__":
    test_weather_server() 