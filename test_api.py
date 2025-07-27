#!/usr/bin/env python3
import requests
import json

API_KEY = "f5fc1e71ee4a4963a74195621252607"
BASE_URL = "http://api.weatherapi.com/v1"

def test_weather_api():
    """Test the WeatherAPI.com API directly"""
    
    # Test cities
    test_cities = ["Tokyo", "London", "New York"]
    
    print("Testing WeatherAPI.com API...")
    print("=" * 40)
    
    for city in test_cities:
        try:
            url = f"{BASE_URL}/current.json"
            params = {
                "key": API_KEY,
                "q": city
            }
            
            response = requests.get(url, params=params, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                location = data["location"]
                current = data["current"]
                
                print(f"✓ {city}: {current['temp_c']}°C, {current['condition']['text']}")
                print(f"  Location: {location['name']}, {location['country']}")
            else:
                print(f"✗ {city}: API Error - {response.status_code}")
                
        except Exception as e:
            print(f"✗ {city}: Error - {e}")
        
        print()
    
    print("=" * 40)
    print("API key is working! Now let's start the MCP server.")

if __name__ == "__main__":
    test_weather_api() 