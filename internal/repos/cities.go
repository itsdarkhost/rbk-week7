package repos

import (
	"context"

	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/jmoiron/sqlx"
)

type CityRepo struct {
	db *sqlx.DB
}

// MARK: New City Repo
func NewCityRepo(db *sqlx.DB) *CityRepo {
	return &CityRepo{db: db}
}

// MARK: Add City
func (r *CityRepo) Create(ctx context.Context, userId int, name string) (*models.City, error) {
	var city models.City

	err := r.db.GetContext(ctx, &city, `
		INSERT INTO cities (user_id, name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, name) DO UPDATE
		SET name = EXCLUDED.name
		RETURNING id, user_id, name
	`, userId, name)

	if err != nil {
		return nil, err
	}

	return &city, nil
}

// MARK: List Cities
func (r *CityRepo) List(ctx context.Context, userId int) ([]models.City, error) {
	var cities []models.City

	err := r.db.SelectContext(ctx, &cities, `
		SELECT id, user_id, name
		FROM cities
		WHERE user_id = $1
		ORDER BY id
	`, userId)

	return cities, err
}

// MARK: Delete City
func (r *CityRepo) Delete(ctx context.Context, userId int, cityId int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM cities
		WHERE user_id = $1
		AND id = $2
	`, userId, cityId)

	return err
}
