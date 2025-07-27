package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WeatherData represents the structure of weather API response
type WeatherData struct {
	Location struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Region  string  `json:"region"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
		Humidity   int     `json:"humidity"`
		WindKph    float64 `json:"wind_kph"`
		WindMph    float64 `json:"wind_mph"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		FeelsLikeC float64 `json:"feelslike_c"`
		FeelsLikeF float64 `json:"feelslike_f"`
		UV         float64 `json:"uv"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempC  float64 `json:"maxtemp_c"`
				MintempC  float64 `json:"mintemp_c"`
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Mark3Labs Weather MCP Server",
		"1.0.0",
	)

	// Add current weather tool with advanced features
	currentWeatherTool := mcp.NewTool("get_current_weather",
		mcp.WithDescription("Get current weather for a city with detailed information"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("City name to get weather for"),
		),
		mcp.WithString("units",
			mcp.Description("Temperature units: celsius or fahrenheit"),
			mcp.Enum("celsius", "fahrenheit"),
		),
	)

	s.AddTool(currentWeatherTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		units, _ := request.OptionalString("units", "celsius")

		// Fetch weather data
		weatherData, err := fetchWeatherData(apiKey, city)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch weather data: %v", err)), nil
		}

		// Format response based on units
		var temp float64
		var tempUnit string
		var feelsLike float64
		if units == "fahrenheit" {
			temp = weatherData.Current.TempF
			tempUnit = "°F"
			feelsLike = weatherData.Current.FeelsLikeF
		} else {
			temp = weatherData.Current.TempC
			tempUnit = "°C"
			feelsLike = weatherData.Current.FeelsLikeC
		}

		result := fmt.Sprintf(`🌤️ Current Weather for %s, %s:

🌡️ Temperature: %.1f%s
🌡️ Feels Like: %.1f%s
☁️ Condition: %s
💧 Humidity: %d%%
💨 Wind: %.1f km/h %s
🌪️ Pressure: %.1f mb
☀️ UV Index: %.1f
📍 Location: %.4f, %.4f

📊 Weather Details:
• Country: %s
• Region: %s
• Wind Speed (mph): %.1f mph
• Last Updated: %s`,
			weatherData.Location.Name,
			weatherData.Location.Country,
			temp,
			tempUnit,
			feelsLike,
			tempUnit,
			weatherData.Current.Condition.Text,
			weatherData.Current.Humidity,
			weatherData.Current.WindKph,
			weatherData.Current.WindDir,
			weatherData.Current.PressureMb,
			weatherData.Current.UV,
			weatherData.Location.Lat,
			weatherData.Location.Lon,
			weatherData.Location.Country,
			weatherData.Location.Region,
			weatherData.Current.WindMph,
			time.Now().Format("2006-01-02 15:04:05"),
		)

		return mcp.NewToolResultText(result), nil
	})

	// Add weather forecast tool with configurable days
	forecastTool := mcp.NewTool("get_weather_forecast",
		mcp.WithDescription("Get weather forecast for a city with configurable number of days"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("City name to get forecast for"),
		),
		mcp.WithInteger("days",
			mcp.Description("Number of days for forecast (1-7)"),
			mcp.Minimum(1),
			mcp.Maximum(7),
		),
		mcp.WithString("units",
			mcp.Description("Temperature units: celsius or fahrenheit"),
			mcp.Enum("celsius", "fahrenheit"),
		),
	)

	s.AddTool(forecastTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		days, _ := request.OptionalInteger("days", 3)
		units, _ := request.OptionalString("units", "celsius")

		// Fetch weather data
		weatherData, err := fetchWeatherData(apiKey, city)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch weather data: %v", err)), nil
		}

		// Format forecast
		result := fmt.Sprintf("🌤️ Weather Forecast for %s, %s (%d days):\n\n", 
			weatherData.Location.Name, weatherData.Location.Country, days)
		
		for i, day := range weatherData.Forecast.Forecastday {
			if i >= days {
				break
			}
			
			var maxTemp, minTemp float64
			var tempUnit string
			if units == "fahrenheit" {
				maxTemp = day.Day.MaxtempC*9/5 + 32
				minTemp = day.Day.MintempC*9/5 + 32
				tempUnit = "°F"
			} else {
				maxTemp = day.Day.MaxtempC
				minTemp = day.Day.MintempC
				tempUnit = "°C"
			}
			
			result += fmt.Sprintf("📅 %s:\n   🌡️ High: %.1f%s\n   🌡️ Low: %.1f%s\n   ☁️ Condition: %s\n\n",
				day.Date,
				maxTemp,
				tempUnit,
				minTemp,
				tempUnit,
				day.Day.Condition.Text,
			)
		}

		return mcp.NewToolResultText(result), nil
	})

	// Add air quality tool
	airQualityTool := mcp.NewTool("get_air_quality",
		mcp.WithDescription("Get air quality information for a city"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("City name to get air quality for"),
		),
	)

	s.AddTool(airQualityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// For demo purposes, return a mock air quality response
		// In a real implementation, you would call an air quality API
		result := fmt.Sprintf(`🌬️ Air Quality for %s:

📊 AQI: 45 (Good)
🌫️ PM2.5: 12 μg/m³
🌫️ PM10: 25 μg/m³
☁️ Ozone: 35 ppb
💨 Nitrogen Dioxide: 18 ppb
🌡️ Carbon Monoxide: 0.8 ppm

🏥 Health Impact: Good air quality. Enjoy outdoor activities.
✅ Recommendations:
   • Safe for outdoor activities
   • Good for sensitive groups
   • No health warnings`,
			city,
		)

		return mcp.NewToolResultText(result), nil
	})

	// Add weather alerts tool
	alertsTool := mcp.NewTool("get_weather_alerts",
		mcp.WithDescription("Get weather alerts and warnings for a city"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("City name to get alerts for"),
		),
	)

	s.AddTool(alertsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Mock weather alerts
		result := fmt.Sprintf(`⚠️ Weather Alerts for %s:

🔴 Severe Weather Warning:
   • Type: Thunderstorm Warning
   • Severity: Moderate
   • Duration: 2 hours
   • Description: Thunderstorms with heavy rain and lightning expected

🟡 Weather Advisory:
   • Type: Wind Advisory
   • Severity: Minor
   • Duration: 4 hours
   • Description: Strong winds up to 25 mph expected

✅ Current Status: No immediate threats
📱 Stay informed with local weather updates`,
			city,
		)

		return mcp.NewToolResultText(result), nil
	})

	// Add resources for historical data
	weatherHistoryResource := mcp.Resource{
		URI:         "weather://history/{city}",
		Name:        "Weather History",
		Description: "Historical weather data for a city",
		MIMEType:    "application/json",
	}

	s.AddResource(weatherHistoryResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract city from URI
		city := "Unknown"
		if args := request.Params.Arguments; args != nil {
			if cityArg, ok := args["city"].(string); ok {
				city = cityArg
			}
		}

		// Mock historical data
		historyData := map[string]interface{}{
			"city": city,
			"history": []map[string]interface{}{
				{
					"date":        time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
					"temperature": 22.5,
					"condition":   "Sunny",
					"humidity":    65,
					"wind_speed":  15.2,
					"pressure":    1013.2,
				},
				{
					"date":        time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
					"temperature": 18.2,
					"condition":   "Cloudy",
					"humidity":    78,
					"wind_speed":  12.8,
					"pressure":    1008.5,
				},
				{
					"date":        time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
					"temperature": 25.1,
					"condition":   "Partly Cloudy",
					"humidity":    58,
					"wind_speed":  18.5,
					"pressure":    1015.8,
				},
			},
		}

		jsonData, _ := json.MarshalIndent(historyData, "", "  ")

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		}, nil
	})

	// Add resource for weather statistics
	weatherStatsResource := mcp.Resource{
		URI:         "weather://stats/{city}",
		Name:        "Weather Statistics",
		Description: "Weather statistics and trends for a city",
		MIMEType:    "application/json",
	}

	s.AddResource(weatherStatsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract city from URI
		city := "Unknown"
		if args := request.Params.Arguments; args != nil {
			if cityArg, ok := args["city"].(string); ok {
				city = cityArg
			}
		}

		// Mock statistics data
		statsData := map[string]interface{}{
			"city": city,
			"statistics": map[string]interface{}{
				"average_temperature": 20.5,
				"max_temperature":     32.1,
				"min_temperature":     8.3,
				"average_humidity":    68.2,
				"rainy_days":          45,
				"sunny_days":          120,
				"windy_days":          23,
			},
			"trends": map[string]interface{}{
				"temperature_trend": "increasing",
				"humidity_trend":    "stable",
				"rainfall_trend":    "decreasing",
			},
		}

		jsonData, _ := json.MarshalIndent(statsData, "", "  ")

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		}, nil
	})

	// Start the server
	log.Println("🚀 Starting Mark3Labs Weather MCP Server...")
	log.Println("📋 Available tools:")
	log.Println("   • get_current_weather - Get current weather with detailed info")
	log.Println("   • get_weather_forecast - Get multi-day weather forecasts")
	log.Println("   • get_air_quality - Get air quality information")
	log.Println("   • get_weather_alerts - Get weather alerts and warnings")
	log.Println("📁 Available resources:")
	log.Println("   • weather://history/{city} - Historical weather data")
	log.Println("   • weather://stats/{city} - Weather statistics and trends")
	log.Println("🔧 Features:")
	log.Println("   • Type-safe parameters with JSON Schema validation")
	log.Println("   • Multiple temperature units (celsius/fahrenheit)")
	log.Println("   • Configurable forecast days (1-7)")
	log.Println("   • Real-time weather data from WeatherAPI.com")
	log.Println("   • Comprehensive error handling")

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func fetchWeatherData(apiKey, city string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=7&aqi=no", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
} 