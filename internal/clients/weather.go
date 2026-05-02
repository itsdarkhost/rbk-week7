package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
		baseURL = "https://wttr.in"
	}

	return &WeatherClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// MARK: Get Weather
func (c *WeatherClient) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	endpoint := fmt.Sprintf("%s/%s?format=j1", c.baseURL, url.PathEscape(city))

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
		return nil, fmt.Errorf("weather api returned status %d", resp.StatusCode)
	}

	var body struct {
		CurrentCondition []struct {
			TempC       string `json:"temp_C"`
			WeatherDesc []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
		} `json:"current_condition"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if len(body.CurrentCondition) == 0 {
		return nil, fmt.Errorf("weather api returned empty response")
	}

	temp, err := strconv.ParseFloat(body.CurrentCondition[0].TempC, 64)
	if err != nil {
		return nil, err
	}

	description := ""
	if len(body.CurrentCondition[0].WeatherDesc) > 0 {
		description = body.CurrentCondition[0].WeatherDesc[0].Value
	}

	return &models.Weather{
		City:        city,
		Temperature: temp,
		Description: description,
	}, nil
}
