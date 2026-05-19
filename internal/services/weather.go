package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/itsdarkhost/rbk-week4/internal/models"
)

type WeatherClient interface {
	GetWeather(ctx context.Context, city string) (*models.Weather, error)
}

type WeatherHistoryRepository interface {
	Create(ctx context.Context, userId int, weather models.Weather) (*models.WeatherHistory, error)
	List(ctx context.Context, userId int, city string, limit int, offset int) ([]models.WeatherHistory, error)
}

type cachedWeather struct {
	weather   models.Weather
	expiresAt time.Time
}

type WeatherService struct {
	userRepo    UserRepository
	cityRepo    CityRepository
	historyRepo WeatherHistoryRepository
	client      WeatherClient
	cache       map[string]cachedWeather
	mu          sync.Mutex
}

// MARK: New Weather Service
func NewWeatherService(userRepo UserRepository, cityRepo CityRepository, historyRepo WeatherHistoryRepository, client WeatherClient) *WeatherService {
	return &WeatherService{
		userRepo:    userRepo,
		cityRepo:    cityRepo,
		historyRepo: historyRepo,
		client:      client,
		cache:       make(map[string]cachedWeather),
	}
}

// MARK: Get Weather
func (s *WeatherService) GetWeather(ctx context.Context, userId int) ([]models.WeatherHistory, error) {
	_, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	cities, err := s.cityRepo.List(ctx, userId)
	if err != nil {
		return nil, err
	}

	type result struct {
		weather *models.Weather
		err     error
	}

	results := make(chan result, len(cities))
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)

		go func(cityName string) {
			defer wg.Done()

			weather, err := s.getCachedWeather(ctx, cityName)
			results <- result{weather: weather, err: err}
		}(city.Name)
	}

	wg.Wait()
	close(results)

	history := make([]models.WeatherHistory, 0, len(cities))
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}

		item, err := s.historyRepo.Create(ctx, userId, *res.weather)
		if err != nil {
			return nil, err
		}

		history = append(history, *item)
	}

	return history, nil
}

// MARK: History
func (s *WeatherService) History(ctx context.Context, userId int, city string, limit int, offset int) (*models.WeatherHistoryResponse, error) {
	_, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	city = strings.TrimSpace(city)
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	history, err := s.historyRepo.List(ctx, userId, city, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.WeatherHistoryResponse{
		UserId:  userId,
		City:    city,
		History: history,
	}, nil
}

func (s *WeatherService) getCachedWeather(ctx context.Context, city string) (*models.Weather, error) {
	key := strings.ToLower(city)
	now := time.Now()

	s.mu.Lock()
	cached, ok := s.cache[key]
	if ok && cached.expiresAt.After(now) {
		weather := cached.weather
		s.mu.Unlock()
		return &weather, nil
	}
	s.mu.Unlock()

	weather, err := s.client.GetWeather(ctx, city)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache[key] = cachedWeather{
		weather:   *weather,
		expiresAt: now.Add(5 * time.Minute),
	}
	s.mu.Unlock()

	return weather, nil
}
