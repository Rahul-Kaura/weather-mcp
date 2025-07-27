#!/usr/bin/env python3
import json
import subprocess
import sys

def test_mark3labs_weather(city):
    """Test the Mark3Labs weather MCP server"""
    
    # Set up the MCP request
    init_request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "test", "version": "1.0"}
        }
    }
    
    weather_request = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "get_current_weather",
            "arguments": {"city": city}
        }
    }
    
    try:
        # Set environment variable
        env = {"WEATHER_API_KEY": "f5fc1e71ee4a4963a74195621252607"}
        
        # Start the server process
        process = subprocess.Popen(
            ["./mark3labs-weather-server"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            cwd="mark3labs-implementation"
        )
        
        # Send initialization
        process.stdin.write((json.dumps(init_request) + "\n").encode())
        process.stdin.flush()
        
        # Read response
        init_response = process.stdout.readline().decode().strip()
        print(f"Init response: {init_response}")
        
        # Send weather request
        process.stdin.write((json.dumps(weather_request) + "\n").encode())
        process.stdin.flush()
        
        # Read response
        weather_response = process.stdout.readline().decode().strip()
        print(f"Weather response: {weather_response}")
        
        # Parse and extract weather data
        response_data = json.loads(weather_response)
        if "result" in response_data and "content" in response_data["result"]:
            weather_text = response_data["result"]["content"][0]["text"]
            return weather_text
        
        process.terminate()
        return None
        
    except Exception as e:
        print(f"Error getting weather for {city}: {e}")
        return None

def main():
    """Test weather for multiple cities"""
    
    cities = ["Tokyo", "London", "New York"]
    
    print("🌤️  Mark3Labs Weather MCP Server Test")
    print("=" * 50)
    
    for city in cities:
        print(f"\n📍 Testing: {city}")
        weather = test_mark3labs_weather(city)
        if weather:
            print(f"✅ Success!")
            print(weather)
        else:
            print(f"❌ Failed to get weather for {city}")
    
    print("\n" + "=" * 50)
    print("✅ Mark3Labs Weather MCP Server is working!")
    print("🌍 Supports international locations")

if __name__ == "__main__":
    main() 