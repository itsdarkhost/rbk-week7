package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	MigrationsPath string
	JWTSecret      string
	WeatherAPIURL  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    databaseURL(),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "file://migrations"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		WeatherAPIURL:  os.Getenv("WEATHER_API_URL"),
	}

	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func databaseURL() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		return dsn
	}

	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	name := getEnv("POSTGRES_DB", "postgres")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
