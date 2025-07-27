package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WeatherAPI response structures
type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Region  string `json:"region"`
}

type Current struct {
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
	FeelslikeC float64 `json:"feelslike_c"`
	FeelslikeF float64 `json:"feelslike_f"`
	Visibility float64 `json:"vis_km"`
	UV         float64 `json:"uv"`
	GustKph    float64 `json:"gust_kph"`
}

type WeatherAPI struct {
	key     string
	baseURL string
	client  *http.Client
}

func NewWeatherAPI(key string) *WeatherAPI {
	return &WeatherAPI{
		key:     key,
		baseURL: "http://api.weatherapi.com/v1",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *WeatherAPI) GetCurrentWeather(ctx context.Context, city string) (*WeatherResponse, error) {
	query := url.Values{
		"key": {w.key},
		"q":   {city},
	}

	request, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		w.baseURL+"/current.json?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := w.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API error: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherResp WeatherResponse
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &weatherResp, nil
}

func formatWeatherResponse(weather *WeatherResponse) string {
	return fmt.Sprintf(`🌤️ Weather for %s, %s

🌡️ Temperature: %.1f°C (%.1f°F)
🌡️ Feels like: %.1f°C (%.1f°F)
☁️ Condition: %s
💧 Humidity: %d%%
💨 Wind: %.1f km/h (%.1f mph) %s
🌪️ Gusts: %.1f km/h
👁️ Visibility: %.1f km
☀️ UV Index: %.1f
🌡️ Pressure: %.1f mb

Last updated: Current conditions`,
		weather.Location.Name,
		weather.Location.Country,
		weather.Current.TempC,
		weather.Current.TempF,
		weather.Current.FeelslikeC,
		weather.Current.FeelslikeF,
		weather.Current.Condition.Text,
		weather.Current.Humidity,
		weather.Current.WindKph,
		weather.Current.WindMph,
		weather.Current.WindDir,
		weather.Current.GustKph,
		weather.Current.Visibility,
		weather.Current.UV,
		weather.Current.PressureMb,
	)
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	// Initialize weather API client
	weatherAPI := NewWeatherAPI(apiKey)

	// Create a new MCP server
	s := server.NewMCPServer(
		"Mark3Labs Weather MCP Server",
		"1.0.0",
		server.WithLogging(),
	)

	// Add current weather tool
	currentWeatherTool := mcp.NewTool("get_current_weather",
		mcp.WithDescription("Get current weather for any city worldwide using WeatherAPI.com"),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("City name to get weather for (e.g., Tokyo, London, New York)"),
		),
	)

	s.AddTool(currentWeatherTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, err := request.RequireString("city")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get weather data from API
		weather, err := weatherAPI.GetCurrentWeather(ctx, city)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get weather for %s: %v", city, err)), nil
		}

		// Format the response
		result := formatWeatherResponse(weather)
		return mcp.NewToolResultText(result), nil
	})

	// Start the server
	log.Println("🚀 Starting Mark3Labs Weather MCP Server...")
	log.Println("🌍 Supports international locations")
	log.Println("🔑 Using WeatherAPI.com")

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
