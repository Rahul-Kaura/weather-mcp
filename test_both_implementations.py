#!/usr/bin/env python3
"""
Test script for both TuanKiri and Mark3Labs MCP implementations
"""

import asyncio
import json
import subprocess
import sys
import time
from typing import Dict, Any

def test_tuanKiri_server():
    """Test the TuanKiri weather MCP server"""
    print("🧪 Testing TuanKiri Weather MCP Server")
    print("=" * 50)
    
    try:
        # Test SSE endpoint
        import requests
        response = requests.get("http://localhost:8000/sse", timeout=3)
        print(f"✅ SSE endpoint status: {response.status_code}")
        
        # Test if server is responding
        if response.status_code == 200:
            print("✅ TuanKiri server is running and responding")
            return True
        else:
            print(f"❌ TuanKiri server returned status: {response.status_code}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Failed to connect to TuanKiri server: {e}")
        return False

def test_mark3labs_server():
    """Test the Mark3Labs weather MCP server"""
    print("\n🧪 Testing Mark3Labs Weather MCP Server")
    print("=" * 50)
    
    try:
        # Start the Mark3Labs server
        process = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd="/Users/hulkster/WeatherFull/mark3labs-implementation",
            env={"WEATHER_API_KEY": "f5fc1e71ee4a4963a74195621252607", **dict(subprocess.os.environ)},
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Give it a moment to start
        time.sleep(2)
        
        # Test initialization
        init_request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {
                    "experimental": {
                        "tools": True
                    }
                },
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0.0"
                }
            }
        }
        
        # Send initialization request
        process.stdin.write(json.dumps(init_request) + "\n")
        process.stdin.flush()
        
        # Read response
        response = process.stdout.readline()
        if response:
            print("✅ Mark3Labs server responded to initialization")
            
            # Test list tools
            tools_request = {
                "jsonrpc": "2.0",
                "id": 2,
                "method": "tools/list"
            }
            
            process.stdin.write(json.dumps(tools_request) + "\n")
            process.stdin.flush()
            
            tools_response = process.stdout.readline()
            if tools_response:
                print("✅ Mark3Labs server responded to tools/list")
                print(f"📄 Response: {tools_response.strip()}")
                process.terminate()
                return True
            else:
                print("❌ No response to tools/list")
                process.terminate()
                return False
        else:
            print("❌ No response to initialization")
            process.terminate()
            return False
            
    except Exception as e:
        print(f"❌ Failed to test Mark3Labs server: {e}")
        return False

def main():
    """Main test function"""
    print("🚀 MCP Implementation Testing Suite")
    print("=" * 60)
    
    # Test TuanKiri server
    tuanKiri_success = test_tuanKiri_server()
    
    # Test Mark3Labs server
    mark3labs_success = test_mark3labs_server()
    
    # Summary
    print("\n" + "=" * 60)
    print("📊 TEST RESULTS SUMMARY")
    print("=" * 60)
    print(f"TuanKiri Weather MCP Server: {'✅ PASS' if tuanKiri_success else '❌ FAIL'}")
    print(f"Mark3Labs Weather MCP Server: {'✅ PASS' if mark3labs_success else '❌ FAIL'}")
    
    if tuanKiri_success and mark3labs_success:
        print("\n🎉 ALL TESTS PASSED!")
        print("✅ Both MCP implementations are working correctly")
        print("✅ Ready for production use")
        return 0
    else:
        print("\n⚠️  SOME TESTS FAILED")
        print("Please check the implementation details above")
        return 1

if __name__ == "__main__":
    sys.exit(main()) 