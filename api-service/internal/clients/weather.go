package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/itsdarkhost/rbk-week4/internal/models"
)

type WeatherClient struct {
	baseURL string
	client  *http.Client
}

// MARK: New Weather Client
func NewWeatherClient(baseURL string) *WeatherClient {
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	return &WeatherClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// MARK: Get Weather
func (c *WeatherClient) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	endpoint := fmt.Sprintf("%s/weather?city=%s", c.baseURL, url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gateway weather api returned status %d", resp.StatusCode)
	}

	var weather models.Weather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}
	if weather.City == "" {
		weather.City = city
	}

	return &weather, nil
}
