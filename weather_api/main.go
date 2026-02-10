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
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient     *redis.Client
	ctx             = context.Background()
	weatherAPIKey   string
	cacheExpiration time.Duration
)

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Conditions  string  `json:"conditions"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
	Timestamp   string  `json:"timestamp"`
}

type VisualCrossingResponse struct {
	ResolvedAddress string `json:"resolvedAddress"`
	CurrentConditions struct {
		Temp       float64 `json:"temp"`
		Conditions string  `json:"conditions"`
		Humidity   float64 `json:"humidity"`
		WindSpeed  float64 `json:"windspeed"`
	} `json:"currentConditions"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	weatherAPIKey = os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	cacheExp, _ := strconv.Atoi(getEnv("CACHE_EXPIRATION", "43200"))
	cacheExpiration = time.Duration(cacheExp) * time.Second

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v. Running without cache.", err)
		redisClient = nil
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/weather", weatherHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Weather API server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if redisClient != nil {
		if cachedData, err := redisClient.Get(ctx, city).Result(); err == nil {
			log.Printf("Cache hit for city: %s", city)
			w.Write([]byte(cachedData))
			return
		}
	}

	log.Printf("Cache miss for city: %s, fetching from API", city)
	weatherData, err := fetchWeatherFromAPI(city)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching weather data: %v", err), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(weatherData)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	if redisClient != nil {
		if err := redisClient.Set(ctx, city, responseJSON, cacheExpiration).Err(); err != nil {
			log.Printf("Error caching data: %v", err)
		}
	}

	w.Write(responseJSON)
}

func fetchWeatherFromAPI(city string) (*WeatherResponse, error) {
	baseURL := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"
	apiURL := fmt.Sprintf("%s/%s?key=%s&unitGroup=metric&include=current",
		baseURL, url.QueryEscape(city), weatherAPIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var vcResponse VisualCrossingResponse
	if err := json.NewDecoder(resp.Body).Decode(&vcResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &WeatherResponse{
		Location:    vcResponse.ResolvedAddress,
		Temperature: vcResponse.CurrentConditions.Temp,
		Conditions:  vcResponse.CurrentConditions.Conditions,
		Humidity:    vcResponse.CurrentConditions.Humidity,
		WindSpeed:   vcResponse.CurrentConditions.WindSpeed,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}
