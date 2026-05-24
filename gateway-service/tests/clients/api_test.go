package clients_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/clients"
)

func TestAPIClientHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Fatalf("expected path /health, got %q", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := clients.NewAPIClient(server.URL)
	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAPIClientHealthStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := clients.NewAPIClient(server.URL)
	if err := client.Health(context.Background()); err == nil {
		t.Fatal("expected status error")
	}
}
