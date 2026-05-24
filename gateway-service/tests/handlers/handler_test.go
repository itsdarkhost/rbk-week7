package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/handlers"
	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/models"
)

type fakeWeatherProvider struct {
	weather *models.Weather
	err     error
	city    string
}

func (f *fakeWeatherProvider) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	f.city = city
	return f.weather, f.err
}

type fakeAPIHealthChecker struct {
	err error
}

func (f fakeAPIHealthChecker) Health(ctx context.Context) error {
	return f.err
}

func TestHealth(t *testing.T) {
	handler := handlers.NewHandler(&fakeWeatherProvider{}, fakeAPIHealthChecker{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestWeather(t *testing.T) {
	provider := &fakeWeatherProvider{
		weather: &models.Weather{City: "Almaty", Temperature: 22, Description: "clear"},
	}
	handler := handlers.NewHandler(provider, fakeAPIHealthChecker{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/weather?city=Almaty", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if provider.city != "Almaty" {
		t.Fatalf("expected city Almaty, got %q", provider.city)
	}
}

func TestWeatherRequiresCity(t *testing.T) {
	handler := handlers.NewHandler(&fakeWeatherProvider{}, fakeAPIHealthChecker{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestAPIHealthReturnsBadGateway(t *testing.T) {
	handler := handlers.NewHandler(&fakeWeatherProvider{}, fakeAPIHealthChecker{err: errors.New("api down")}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api-health", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}
