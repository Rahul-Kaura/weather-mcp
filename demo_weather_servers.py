#!/usr/bin/env python3
"""
Demo script showing the capabilities of all three weather MCP servers
"""

import json
import subprocess
import time

def demo_python_server():
    """Demo the Python weather server"""
    print("🐍 Python Weather MCP Server Demo")
    print("=" * 40)
    print("Available tools:")
    print("• get_alerts(state) - Get weather alerts for US states")
    print("• get_forecast(lat, lon) - Get weather forecast for coordinates")
    print("• Uses National Weather Service API (free)")
    print("• STDIO transport")
    print()

def demo_tuanKiri_server():
    """Demo the TuanKiri weather server"""
    print("🚀 TuanKiri Weather MCP Server Demo")
    print("=" * 40)
    print("Available tools:")
    print("• current_weather(city) - Get current weather for any city")
    print("• Uses WeatherAPI.com for real-time data")
    print("• HTTP/SSE transport on localhost:8000")
    print("• Simple and fast")
    print()

def demo_mark3labs_server():
    """Demo the Mark3Labs weather server"""
    print("⚡ Mark3Labs Weather MCP Server Demo")
    print("=" * 40)
    print("Available tools:")
    print("• get_current_weather(city, units) - Detailed current weather")
    print("• get_weather_forecast(city, days, units) - Multi-day forecasts")
    print("• get_air_quality(city) - Air quality information")
    print("• get_weather_alerts(city) - Weather alerts and warnings")
    print("• Uses WeatherAPI.com with advanced formatting")
    print("• STDIO transport with type-safe parameters")
    print("• Resources: weather://history/{city}, weather://stats/{city}")
    print()

def show_claude_integration():
    """Show how to use with Claude"""
    print("🤖 Claude for Desktop Integration")
    print("=" * 40)
    print("All three servers are now configured in Claude!")
    print()
    print("To use with Claude:")
    print("1. Restart Claude for Desktop")
    print("2. Ask Claude: 'What's the weather in Tokyo?'")
    print("3. Claude will automatically use the weather tools")
    print()
    print("Example conversations:")
    print("• 'Get weather alerts for California'")
    print("• 'What's the 5-day forecast for New York?'")
    print("• 'Check air quality in London'")
    print("• 'Get weather statistics for Paris'")
    print()

def main():
    """Main demo function"""
    print("🌤️ Weather MCP Servers - Complete Demo")
    print("=" * 60)
    print()
    
    demo_python_server()
    demo_tuanKiri_server()
    demo_mark3labs_server()
    show_claude_integration()
    
    print("🎉 All three implementations are ready to use!")
    print("✅ Python server: National Weather Service data")
    print("✅ TuanKiri server: Simple, fast weather data")
    print("✅ Mark3Labs server: Advanced features and formatting")
    print()
    print("🚀 Next: Restart Claude for Desktop and start asking about weather!")

if __name__ == "__main__":
    main() 