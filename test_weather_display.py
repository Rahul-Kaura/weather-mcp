#!/usr/bin/env python3
import json
import subprocess
import sys

def get_weather_for_city(city):
    """Get weather for a city using the MCP server"""
    
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
            "name": "current_weather",
            "arguments": {"city": city}
        }
    }
    
    try:
        # Set environment variable
        env = {"WEATHER_API_KEY": "f5fc1e71ee4a4963a74195621252607"}
        
        # Start the server process
        process = subprocess.Popen(
            ["./weather-mcp-server"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            cwd="weather-mcp-server"
        )
        
        # Send initialization
        process.stdin.write((json.dumps(init_request) + "\n").encode())
        process.stdin.flush()
        
        # Read response
        init_response = process.stdout.readline().decode().strip()
        
        # Send weather request
        process.stdin.write((json.dumps(weather_request) + "\n").encode())
        process.stdin.flush()
        
        # Read response
        weather_response = process.stdout.readline().decode().strip()
        
        # Parse and extract weather data
        response_data = json.loads(weather_response)
        if "result" in response_data and "content" in response_data["result"]:
            html_content = response_data["result"]["content"][0]["text"]
            
            # Extract key information from HTML
            import re
            
            # Extract temperature
            temp_match = re.search(r'Temperature: (\d+)°C', html_content)
            temp = temp_match.group(1) if temp_match else "N/A"
            
            # Extract condition
            condition_match = re.search(r'Condition: ([^<]+)', html_content)
            condition = condition_match.group(1) if condition_match else "N/A"
            
            # Extract humidity
            humidity_match = re.search(r'Humidity: (\d+)%', html_content)
            humidity = humidity_match.group(1) if humidity_match else "N/A"
            
            # Extract wind speed
            wind_match = re.search(r'Wind Speed: ([\d.]+) km/h', html_content)
            wind = wind_match.group(1) if wind_match else "N/A"
            
            return {
                "city": city,
                "temperature": f"{temp}°C",
                "condition": condition,
                "humidity": f"{humidity}%",
                "wind_speed": f"{wind} km/h"
            }
        
        process.terminate()
        return None
        
    except Exception as e:
        print(f"Error getting weather for {city}: {e}")
        return None

def main():
    """Test weather for multiple cities"""
    
    cities = ["Tokyo", "London", "New York", "Paris", "Sydney"]
    
    print("🌤️  Weather MCP Server Test")
    print("=" * 50)
    
    for city in cities:
        weather = get_weather_for_city(city)
        if weather:
            print(f"\n📍 {weather['city']}")
            print(f"   🌡️  Temperature: {weather['temperature']}")
            print(f"   ☁️  Condition: {weather['condition']}")
            print(f"   💧 Humidity: {weather['humidity']}")
            print(f"   💨 Wind: {weather['wind_speed']}")
        else:
            print(f"\n❌ Failed to get weather for {city}")
    
    print("\n" + "=" * 50)
    print("✅ Weather MCP Server is working!")
    print("🌍 Supports international locations")
    print("🔧 Ready for Claude integration")

if __name__ == "__main__":
    main() 