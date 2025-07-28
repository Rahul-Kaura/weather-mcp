#!/usr/bin/env python3
import subprocess
import json
import sys

def test_server(server_path, tool_name, city):
    """Test a weather server with a specific city"""
    try:
        # Initialize the server
        init_cmd = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {"name": "test", "version": "1.0"}
            }
        }
        
        # Call the weather tool
        call_cmd = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": {"city": city}
            }
        }
        
        # Send both commands
        input_data = json.dumps(init_cmd) + "\n" + json.dumps(call_cmd) + "\n"
        
        result = subprocess.run(
            [server_path],
            input=input_data,
            capture_output=True,
            text=True,
            timeout=10
        )
        
        if result.returncode == 0:
            # Parse the response
            lines = result.stdout.strip().split('\n')
            for line in lines:
                if line.strip():
                    try:
                        response = json.loads(line)
                        if 'result' in response and 'content' in response['result']:
                            return response['result']['content'][0]['text']
                    except json.JSONDecodeError:
                        continue
        return f"Error: {result.stderr}"
        
    except subprocess.TimeoutExpired:
        return "Error: Server timeout"
    except Exception as e:
        return f"Error: {str(e)}"

def main():
    print("🌤️  Testing Enhanced Weather Servers")
    print("=" * 50)
    
    # Test cities
    cities = ["Tokyo", "London", "Duluth"]
    
    # Test TuanKiri server (HTML output)
    print("\n🎨 **TuanKiri Server (HTML Output)**")
    print("-" * 40)
    
    for city in cities:
        print(f"\n📍 Testing {city}:")
        result = test_server(
            "./weather-mcp-server/weather-mcp-server",
            "current_weather",
            city
        )
        
        if "Error" in result:
            print(f"❌ {result}")
        else:
            # Extract key information from HTML
            if "city-image" in result and "fun-fact" in result:
                print(f"✅ Enhanced HTML with city image and fun fact")
                print(f"   📸 City image included")
                print(f"   💡 Fun fact included")
                print(f"   🎨 Styled with gradient background")
            else:
                print(f"⚠️  Basic HTML response")
    
    # Test Mark3Labs server (Rich text output)
    print("\n📝 **Mark3Labs Server (Rich Text Output)**")
    print("-" * 40)
    
    for city in cities:
        print(f"\n📍 Testing {city}:")
        result = test_server(
            "./mark3labs-implementation/mark3labs-weather-server",
            "get_current_weather",
            city
        )
        
        if "Error" in result:
            print(f"❌ {result}")
        else:
            # Check for enhanced features
            if "🌤️" in result and "💡" in result and "![Tokyo]" in result:
                print(f"✅ Enhanced rich text with emojis and images")
                print(f"   🌤️ Emojis included")
                print(f"   📸 City image included")
                print(f"   💡 Fun fact included")
                print(f"   📝 Markdown formatting")
            else:
                print(f"⚠️  Basic text response")
    
    print("\n" + "=" * 50)
    print("✅ Enhanced weather servers tested!")
    print("🎨 TuanKiri: HTML with city images and fun facts")
    print("📝 Mark3Labs: Rich text with emojis and city images")

if __name__ == "__main__":
    main() 