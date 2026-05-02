package repos

import (
	"context"

	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/jmoiron/sqlx"
)

type WeatherHistoryRepo struct {
	db *sqlx.DB
}

// MARK: New Weather History Repo
func NewWeatherHistoryRepo(db *sqlx.DB) *WeatherHistoryRepo {
	return &WeatherHistoryRepo{db: db}
}

// MARK: Create History
func (r *WeatherHistoryRepo) Create(ctx context.Context, userId int, weather models.Weather) (*models.WeatherHistory, error) {
	var history models.WeatherHistory

	err := r.db.GetContext(ctx, &history, `
		INSERT INTO weather_history (user_id, city, temperature, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, city, temperature, description, requested_at
	`, userId, weather.City, weather.Temperature, weather.Description)

	if err != nil {
		return nil, err
	}

	return &history, nil
}

// MARK: List History
func (r *WeatherHistoryRepo) List(ctx context.Context, userId int, city string, limit int, offset int) ([]models.WeatherHistory, error) {
	var history []models.WeatherHistory

	if city != "" {
		err := r.db.SelectContext(ctx, &history, `
			SELECT id, user_id, city, temperature, description, requested_at
			FROM weather_history
			WHERE user_id = $1
			AND city = $2
			ORDER BY requested_at DESC
			LIMIT $3 OFFSET $4
		`, userId, city, limit, offset)

		return history, err
	}

	err := r.db.SelectContext(ctx, &history, `
		SELECT id, user_id, city, temperature, description, requested_at
		FROM weather_history
		WHERE user_id = $1
		ORDER BY requested_at DESC
		LIMIT $2 OFFSET $3
	`, userId, limit, offset)

	return history, err
}
