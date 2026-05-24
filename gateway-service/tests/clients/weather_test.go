package clients_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/clients"
)

func TestWeatherClientGetWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Almaty" {
			t.Fatalf("expected path /Almaty, got %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"current_condition": [
				{
					"temp_C": "21.5",
					"weatherDesc": [{"value": "Sunny"}]
				}
			]
		}`))
	}))
	defer server.Close()

	client := clients.NewWeatherClient(server.URL)
	weather, err := client.GetWeather(context.Background(), "Almaty")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if weather.City != "Almaty" {
		t.Fatalf("expected city Almaty, got %q", weather.City)
	}
	if weather.Temperature != 21.5 {
		t.Fatalf("expected temperature 21.5, got %v", weather.Temperature)
	}
	if weather.Description != "Sunny" {
		t.Fatalf("expected description Sunny, got %q", weather.Description)
	}
}
