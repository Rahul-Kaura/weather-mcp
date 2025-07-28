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
	"strings"
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

// getCityImage returns a weather-dependent location recommendation based on the city name
func getCityImage(city string, condition string, tempC float64) string {
	cityLower := strings.ToLower(city)

	// Weather-dependent location recommendations
	locationRecommendations := map[string]map[string]string{
		"tokyo": {
			"hot":   "🏯 Visit the air-conditioned Tokyo National Museum",
			"warm":  "🌸 Stroll through Shinjuku Gyoen National Garden",
			"cool":  "🗼 Climb Tokyo Tower for city views",
			"cold":  "♨️ Relax in a traditional onsen (hot spring)",
			"rainy": "🏛️ Explore the Imperial Palace East Gardens",
			"sunny": "🎌 Walk the historic Meiji Shrine",
		},
		"london": {
			"hot":   "🏛️ Cool off at the British Museum",
			"warm":  "🌳 Enjoy Hyde Park and Kensington Gardens",
			"cool":  "🎭 Visit the West End theatres",
			"cold":  "☕ Warm up in a traditional English pub",
			"rainy": "🏛️ Explore the Natural History Museum",
			"sunny": "🌉 Walk across Tower Bridge",
		},
		"new york": {
			"hot":   "🏛️ Visit the air-conditioned Metropolitan Museum",
			"warm":  "🌳 Stroll through Central Park",
			"cool":  "🗽 Take the ferry to Statue of Liberty",
			"cold":  "☕ Warm up in a cozy Brooklyn café",
			"rainy": "🎭 Catch a Broadway show",
			"sunny": "🌆 Walk the High Line elevated park",
		},
		"paris": {
			"hot":   "🏛️ Cool off at the Louvre Museum",
			"warm":  "🌸 Stroll through Luxembourg Gardens",
			"cool":  "🗼 Visit the Eiffel Tower",
			"cold":  "☕ Warm up in a charming café",
			"rainy": "🏛️ Explore the Musée d'Orsay",
			"sunny": "🌉 Walk along the Seine River",
		},
		"sydney": {
			"hot":   "🏛️ Visit the air-conditioned Art Gallery of NSW",
			"warm":  "🏖️ Relax at Bondi Beach",
			"cool":  "🎭 Visit the Sydney Opera House",
			"cold":  "☕ Warm up in a harbor-side café",
			"rainy": "🏛️ Explore the Australian Museum",
			"sunny": "🌉 Walk across Sydney Harbour Bridge",
		},
		"duluth": {
			"hot":   "🏛️ Visit the Great Lakes Aquarium",
			"warm":  "🌊 Walk along Lake Superior",
			"cool":  "🌉 Visit the Aerial Lift Bridge",
			"cold":  "☕ Warm up in a cozy café",
			"rainy": "🏛️ Explore the Duluth Art Institute",
			"sunny": "🌳 Hike in Enger Park",
		},
		"mumbai": {
			"hot":   "🏛️ Visit the air-conditioned National Museum",
			"warm":  "🌊 Walk along Marine Drive",
			"cool":  "🏛️ Visit the Gateway of India",
			"cold":  "☕ Warm up in a local café",
			"rainy": "🏛️ Explore the Chhatrapati Shivaji Museum",
			"sunny": "🌳 Visit the Sanjay Gandhi National Park",
		},
		"beijing": {
			"hot":   "🏛️ Visit the air-conditioned National Museum",
			"warm":  "🏯 Walk through the Forbidden City",
			"cool":  "🐉 Visit the Temple of Heaven",
			"cold":  "☕ Warm up in a traditional tea house",
			"rainy": "🏛️ Explore the Capital Museum",
			"sunny": "🌉 Walk along the Great Wall",
		},
		"moscow": {
			"hot":   "🏛️ Visit the air-conditioned Tretyakov Gallery",
			"warm":  "🌳 Stroll through Gorky Park",
			"cool":  "⛪ Visit Saint Basil's Cathedral",
			"cold":  "☕ Warm up in a cozy café",
			"rainy": "🏛️ Explore the Pushkin Museum",
			"sunny": "🏰 Walk through Red Square",
		},
		"cairo": {
			"hot":   "🏛️ Visit the air-conditioned Egyptian Museum",
			"warm":  "🐪 Take a camel ride near the pyramids",
			"cool":  "🏺 Visit the Great Pyramid of Giza",
			"cold":  "☕ Warm up in a traditional café",
			"rainy": "🏛️ Explore the Coptic Museum",
			"sunny": "🌊 Take a Nile River cruise",
		},
	}

	// Determine weather category
	var weatherCategory string
	switch {
	case tempC >= 25:
		weatherCategory = "hot"
	case tempC >= 15:
		weatherCategory = "warm"
	case tempC >= 5:
		weatherCategory = "cool"
	default:
		weatherCategory = "cold"
	}

	// Check for specific weather conditions
	if strings.Contains(strings.ToLower(condition), "rain") {
		weatherCategory = "rainy"
	} else if strings.Contains(strings.ToLower(condition), "sunny") || strings.Contains(strings.ToLower(condition), "clear") {
		weatherCategory = "sunny"
	}

	// Get recommendation for the city
	if cityRecs, exists := locationRecommendations[cityLower]; exists {
		if rec, exists := cityRecs[weatherCategory]; exists {
			return rec
		}
		// Fallback to warm if specific category doesn't exist
		if rec, exists := cityRecs["warm"]; exists {
			return rec
		}
	}

	// Default recommendations for other cities
	defaultRecs := map[string]string{
		"hot":   "🏛️ Visit a local museum to cool off",
		"warm":  "🌳 Enjoy a local park or garden",
		"cool":  "🏛️ Explore local attractions",
		"cold":  "☕ Warm up in a cozy café",
		"rainy": "🏛️ Visit indoor attractions",
		"sunny": "🌳 Enjoy outdoor activities",
	}

	return defaultRecs[weatherCategory]
}

// getFunFact returns a weather-dependent fun fact about the city
func getFunFact(city string, condition string, tempC float64) string {
	cityLower := strings.ToLower(city)
	conditionLower := strings.ToLower(condition)

	// Weather-dependent fun facts
	if strings.Contains(conditionLower, "rain") || strings.Contains(conditionLower, "drizzle") {
		switch cityLower {
		case "tokyo":
			return "Perfect weather for visiting the beautiful cherry blossoms in Ueno Park! 🌸"
		case "london":
			return "Classic London weather! Great time to visit the British Museum or enjoy a cozy pub. ☔"
		case "paris":
			return "Rainy days in Paris are perfect for exploring the Louvre or enjoying café culture! ☕"
		case "sydney":
			return "Rain in Sydney means the Opera House looks even more dramatic! 🎭"
		default:
			return "Rainy days are perfect for indoor activities and cozy cafes! ☔"
		}
	} else if strings.Contains(conditionLower, "sunny") || strings.Contains(conditionLower, "clear") {
		switch cityLower {
		case "tokyo":
			return "Perfect weather for visiting Tokyo Tower or taking a stroll in Yoyogi Park! 🗼"
		case "london":
			return "Sunny London! Great time to visit Hyde Park or take a Thames River cruise! ☀️"
		case "new york":
			return "Perfect weather for walking across the Brooklyn Bridge or visiting Central Park! 🌉"
		case "paris":
			return "Beautiful weather for climbing the Eiffel Tower or strolling along the Seine! 🗼"
		case "sydney":
			return "Great weather for visiting Bondi Beach or climbing the Sydney Harbour Bridge! 🏖️"
		default:
			return "Perfect weather for outdoor activities and sightseeing! ☀️"
		}
	} else if strings.Contains(conditionLower, "cloudy") || strings.Contains(conditionLower, "overcast") {
		switch cityLower {
		case "duluth":
			return "Cloudy weather is perfect for exploring the beautiful Lake Superior shoreline! 🏞️"
		case "london":
			return "Classic London weather! Great for visiting museums or enjoying afternoon tea! ☁️"
		case "sydney":
			return "Cloudy weather is perfect for visiting the Royal Botanic Garden! 🌿"
		default:
			return "Cloudy weather is great for photography and exploring indoor attractions! ☁️"
		}
	} else if tempC > 30 {
		return "Hot weather! Perfect time for ice cream and finding air-conditioned spots! 🍦"
	} else if tempC < 10 {
		return "Cold weather! Great time for hot drinks and cozy indoor activities! ☕"
	}

	// Default fun facts by city
	switch cityLower {
	case "tokyo":
		return "Did you know Tokyo has the world's busiest pedestrian crossing at Shibuya? 🚶‍♂️"
	case "london":
		return "London has over 170 museums, many of them free to visit! 🏛️"
	case "new york":
		return "New York City has over 8 million people and 800 languages spoken! 🌆"
	case "paris":
		return "Paris is known as the 'City of Light' and has over 300 illuminated monuments! 💡"
	case "sydney":
		return "Sydney Harbour is home to over 600 species of fish! 🐟"
	case "duluth":
		return "Duluth is home to the world's largest freshwater lake, Lake Superior! 🏞️"
	default:
		return "Every city has its unique charm and hidden gems to discover! ✨"
	}
}

// getWeatherAlert returns weather alerts based on conditions
func getWeatherAlert(condition string, tempC float64, windKph float64, humidity int) string {
	conditionLower := strings.ToLower(condition)

	if strings.Contains(conditionLower, "thunder") {
		return "⚡ **WEATHER ALERT:** Thunderstorm detected - seek shelter immediately!"
	} else if strings.Contains(conditionLower, "storm") {
		return "🌪️ **WEATHER ALERT:** Storm conditions - avoid outdoor activities!"
	} else if tempC > 35 {
		return "🔥 **HEAT ALERT:** Extreme heat - stay hydrated and avoid sun exposure!"
	} else if tempC < -10 {
		return "❄️ **COLD ALERT:** Extreme cold - bundle up and limit outdoor time!"
	} else if windKph > 50 {
		return "💨 **WIND ALERT:** High winds - secure loose items and be cautious!"
	} else if humidity > 90 {
		return "💧 **HUMIDITY ALERT:** Very high humidity - stay hydrated!"
	}

	return ""
}

// getAirQualityRecommendation returns air quality advice
func getAirQualityRecommendation(uv float64, visibility float64) string {
	if uv > 8 {
		return "☀️ **UV Index High:** Use SPF 50+ sunscreen and limit sun exposure"
	} else if uv > 6 {
		return "☀️ **UV Index Moderate:** Use sunscreen and wear protective clothing"
	} else if visibility < 5 {
		return "🌫️ **Poor Visibility:** Drive carefully and avoid outdoor activities"
	}

	return "✅ **Good Conditions:** Safe for outdoor activities"
}

// getTravelRecommendations returns travel-specific advice
func getTravelRecommendations(city string, condition string, tempC float64) string {
	cityLower := strings.ToLower(city)
	conditionLower := strings.ToLower(condition)

	var recommendations []string

	// General travel advice
	if strings.Contains(conditionLower, "rain") {
		recommendations = append(recommendations, "☔ **Travel Tip:** Pack waterproof gear and plan indoor activities")
	} else if strings.Contains(conditionLower, "sunny") {
		recommendations = append(recommendations, "🧴 **Travel Tip:** Bring sunscreen and stay hydrated")
	} else if tempC > 30 {
		recommendations = append(recommendations, "🏊 **Travel Tip:** Perfect weather for water activities")
	} else if tempC < 10 {
		recommendations = append(recommendations, "🧥 **Travel Tip:** Pack warm clothing and plan indoor visits")
	}

	// City-specific travel advice
	switch cityLower {
	case "tokyo":
		recommendations = append(recommendations, "🚇 **Local Tip:** Use the efficient subway system to get around")
		recommendations = append(recommendations, "🍜 **Food Tip:** Try local ramen shops for authentic cuisine")
	case "london":
		recommendations = append(recommendations, "🚇 **Local Tip:** Get an Oyster card for public transport")
		recommendations = append(recommendations, "☕ **Food Tip:** Experience traditional afternoon tea")
	case "new york":
		recommendations = append(recommendations, "🚇 **Local Tip:** Use the subway - it's the fastest way around")
		recommendations = append(recommendations, "🍕 **Food Tip:** Try authentic New York pizza")
	case "paris":
		recommendations = append(recommendations, "🚇 **Local Tip:** Use the Metro for easy navigation")
		recommendations = append(recommendations, "🥐 **Food Tip:** Visit local bakeries for fresh pastries")
	case "sydney":
		recommendations = append(recommendations, "🚇 **Local Tip:** Use Opal card for public transport")
		recommendations = append(recommendations, "🦘 **Local Tip:** Visit wildlife parks to see native animals")
	case "duluth":
		recommendations = append(recommendations, "🚗 **Local Tip:** Rent a car to explore the scenic shoreline")
		recommendations = append(recommendations, "🏞️ **Local Tip:** Visit state parks for hiking and nature")
	}

	return strings.Join(recommendations, "\n")
}

// getWeatherScore returns a weather score out of 10
func getWeatherScore(condition string, tempC float64, humidity int, windKph float64) (int, string) {
	score := 5 // Base score

	// Temperature scoring
	if tempC >= 18 && tempC <= 25 {
		score += 3 // Perfect temperature
	} else if tempC >= 15 && tempC <= 28 {
		score += 2 // Good temperature
	} else if tempC >= 10 && tempC <= 32 {
		score += 1 // Acceptable temperature
	}

	// Condition scoring
	conditionLower := strings.ToLower(condition)
	if strings.Contains(conditionLower, "sunny") || strings.Contains(conditionLower, "clear") {
		score += 2
	} else if strings.Contains(conditionLower, "partly cloudy") {
		score += 1
	} else if strings.Contains(conditionLower, "rain") {
		score -= 1
	} else if strings.Contains(conditionLower, "storm") {
		score -= 2
	}

	// Wind scoring
	if windKph < 20 {
		score += 1
	} else if windKph > 40 {
		score -= 1
	}

	// Humidity scoring
	if humidity >= 40 && humidity <= 70 {
		score += 1
	} else if humidity > 80 {
		score -= 1
	}

	// Clamp score between 1 and 10
	if score > 10 {
		score = 10
	} else if score < 1 {
		score = 1
	}

	// Generate score description
	var description string
	switch {
	case score >= 9:
		description = "🌟 **Excellent weather conditions!**"
	case score >= 7:
		description = "👍 **Great weather for activities!**"
	case score >= 5:
		description = "😊 **Decent weather conditions**"
	case score >= 3:
		description = "😐 **Moderate weather - plan accordingly**"
	default:
		description = "😔 **Challenging weather conditions**"
	}

	return score, description
}

func formatWeatherResponse(weather *WeatherResponse, city string) string {
	cityImage := getCityImage(city, weather.Current.Condition.Text, weather.Current.TempC)
	funFact := getFunFact(city, weather.Current.Condition.Text, weather.Current.TempC)
	weatherAlert := getWeatherAlert(weather.Current.Condition.Text, weather.Current.TempC, weather.Current.WindKph, weather.Current.Humidity)
	airQuality := getAirQualityRecommendation(weather.Current.UV, weather.Current.Visibility)
	travelTips := getTravelRecommendations(city, weather.Current.Condition.Text, weather.Current.TempC)
	weatherScore, scoreDescription := getWeatherScore(weather.Current.Condition.Text, weather.Current.TempC, weather.Current.Humidity, weather.Current.WindKph)

	// Build the enhanced response
	response := fmt.Sprintf(`🌤️ **Weather for %s, %s** 🌤️

📍 **Location Recommendation:**
%s

**Current Conditions:**
🌡️ **Temperature:** %.1f°C (%.1f°F)
🌡️ **Feels like:** %.1f°C (%.1f°F)
☁️ **Condition:** %s
💧 **Humidity:** %d%%
💨 **Wind:** %.1f km/h (%.1f mph) %s
🌪️ **Gusts:** %.1f km/h
👁️ **Visibility:** %.1f km
☀️ **UV Index:** %.1f
🌡️ **Pressure:** %.1f mb

**Weather Score:** %d/10 %s

**Weather Icon:** ![Weather](https://www.weatherapi.com/static/img/weather/%s)

💡 **Fun Fact:** %s

%s

%s

**🎯 Travel Recommendations:**
%s

*Last updated: Current conditions*`,
		weather.Location.Name,
		weather.Location.Country,
		weather.Location.Name,
		cityImage,
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
		weatherScore,
		scoreDescription,
		weather.Current.Condition.Icon,
		funFact,
		weatherAlert,
		airQuality,
		travelTips,
	)

	return response
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
		mcp.WithDescription("Get current weather for any city worldwide using WeatherAPI.com with enhanced formatting and fun facts"),
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

		// Format the response with enhanced features
		result := formatWeatherResponse(weather, city)
		return mcp.NewToolResultText(result), nil
	})

	// Start the server
	log.Println("🚀 Starting Mark3Labs Weather MCP Server...")
	log.Println("🌍 Supports international locations")
	log.Println("🔑 Using WeatherAPI.com")
	log.Println("✨ Enhanced with city images and fun facts")

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
