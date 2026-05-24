package config

import "os"

type Config struct {
	Port                  string
	ExternalWeatherAPIURL string
	APIServiceURL         string
}

func Load() Config {
	return Config{
		Port:                  getEnv("PORT", "8081"),
		ExternalWeatherAPIURL: getEnv("EXTERNAL_WEATHER_API_URL", "https://wttr.in"),
		APIServiceURL:         getEnv("API_SERVICE_URL", "http://localhost:8080"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
