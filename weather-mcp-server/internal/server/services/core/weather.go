package core

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type WeatherService struct {
	*CoreServices
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

// getWeatherRecommendations returns personalized recommendations based on weather
func getWeatherRecommendations(city string, condition string, tempC float64, humidity int, windKph float64) string {
	cityLower := strings.ToLower(city)
	conditionLower := strings.ToLower(condition)

	var recommendations []string

	// Temperature-based recommendations
	if tempC > 30 {
		recommendations = append(recommendations, "🌡️ Stay hydrated and avoid prolonged sun exposure")
		recommendations = append(recommendations, "🏊 Perfect weather for swimming or water activities")
	} else if tempC > 20 {
		recommendations = append(recommendations, "☀️ Ideal temperature for outdoor activities")
		recommendations = append(recommendations, "🚶 Great for walking tours and sightseeing")
	} else if tempC > 10 {
		recommendations = append(recommendations, "🧥 Light jacket recommended")
		recommendations = append(recommendations, "☕ Perfect for café visits and indoor activities")
	} else {
		recommendations = append(recommendations, "🧣 Bundle up! Warm clothing essential")
		recommendations = append(recommendations, "🔥 Great time for hot drinks and cozy indoor spots")
	}

	// Condition-based recommendations
	if strings.Contains(conditionLower, "rain") {
		recommendations = append(recommendations, "☔ Bring an umbrella or raincoat")
		recommendations = append(recommendations, "🏛️ Perfect for museum visits and indoor attractions")
	} else if strings.Contains(conditionLower, "sunny") {
		recommendations = append(recommendations, "🧴 Don't forget sunscreen!")
		recommendations = append(recommendations, "📸 Excellent conditions for photography")
	} else if strings.Contains(conditionLower, "cloudy") {
		recommendations = append(recommendations, "📷 Great lighting for photography")
		recommendations = append(recommendations, "🚶 Comfortable for outdoor activities")
	}

	// Wind-based recommendations
	if windKph > 30 {
		recommendations = append(recommendations, "💨 Strong winds - secure loose items")
		recommendations = append(recommendations, "🏠 Consider indoor activities")
	}

	// Humidity-based recommendations
	if humidity > 80 {
		recommendations = append(recommendations, "💧 High humidity - stay hydrated")
		recommendations = append(recommendations, "🌬️ Seek air-conditioned spaces")
	}

	// City-specific recommendations
	switch cityLower {
	case "tokyo":
		recommendations = append(recommendations, "🗼 Visit Tokyo Tower for amazing city views")
		recommendations = append(recommendations, "🌸 Check out local parks and gardens")
	case "london":
		recommendations = append(recommendations, "🏛️ Explore the British Museum")
		recommendations = append(recommendations, "☕ Enjoy traditional afternoon tea")
	case "new york":
		recommendations = append(recommendations, "🌉 Walk across the Brooklyn Bridge")
		recommendations = append(recommendations, "🏙️ Visit Central Park for nature")
	case "paris":
		recommendations = append(recommendations, "🗼 Climb the Eiffel Tower")
		recommendations = append(recommendations, "☕ Experience café culture")
	case "sydney":
		recommendations = append(recommendations, "🏖️ Visit Bondi Beach")
		recommendations = append(recommendations, "🎭 Explore the Opera House")
	case "duluth":
		recommendations = append(recommendations, "🏞️ Explore Lake Superior shoreline")
		recommendations = append(recommendations, "🚢 Visit the maritime museum")
	}

	return strings.Join(recommendations, "\n")
}

// getWeatherTrend returns a simple trend analysis
func getWeatherTrend(tempC float64, condition string) string {
	if tempC > 25 && strings.Contains(strings.ToLower(condition), "sunny") {
		return "📈 Perfect summer weather - great for outdoor activities!"
	} else if tempC < 10 && strings.Contains(strings.ToLower(condition), "cloudy") {
		return "📉 Cool and overcast - indoor activities recommended"
	} else if strings.Contains(strings.ToLower(condition), "rain") {
		return "🌧️ Rainy conditions - bring protection and plan indoor activities"
	} else {
		return "🌤️ Moderate conditions - suitable for most activities"
	}
}

func (ws *WeatherService) Current(ctx context.Context, city string) (string, error) {
	data, err := ws.weatherAPI.Current(ctx, city)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	recommendations := getWeatherRecommendations(city, data.Current.Condition.Text, data.Current.TempC, int(data.Current.Humidity), data.Current.WindKph)
	weatherTrend := getWeatherTrend(data.Current.TempC, data.Current.Condition.Text)

	// Split recommendations into a list for the template
	recommendationsList := strings.Split(recommendations, "\n")

	if err := ws.renderer.ExecuteTemplate(&buf, "weather.html", map[string]interface{}{
		"Location":            fmt.Sprintf("%s, %s", data.Location.Name, data.Location.Country),
		"Icon":                "https:" + data.Current.Condition.Icon,
		"Condition":           data.Current.Condition.Text,
		"Temperature":         fmt.Sprintf("%.0f", data.Current.TempC),
		"Humidity":            fmt.Sprintf("%d", data.Current.Humidity),
		"WindSpeed":           fmt.Sprintf("%.0f", data.Current.WindKph),
		"FeelsLike":           fmt.Sprintf("%.0f", data.Current.FeelslikeC),
		"CityImage":           getCityImage(city, data.Current.Condition.Text, data.Current.TempC),
		"FunFact":             getFunFact(city, data.Current.Condition.Text, data.Current.TempC),
		"WeatherTrend":        weatherTrend,
		"RecommendationsList": recommendationsList,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
