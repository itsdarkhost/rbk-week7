package repos

import (
	"context"

	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

// MARK: New Repo
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// MARK: List
func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.SelectContext(ctx, &users, `
		SELECT id, username, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`)

	return users, err
}

// MARK: Get
func (r *UserRepo) Get(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `
		SELECT id, username, deleted_at
		FROM users
		WHERE id = $1
		AND deleted_at IS NULL
	`, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// MARK: Create
func (r *UserRepo) Create(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `
		INSERT INTO users (username)
		VALUES ($1)
		RETURNING id, username, deleted_at
	`, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// MARK: Update
func (r *UserRepo) Update(ctx context.Context, id int, username string) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `
		UPDATE users
		SET username = $1
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING id, username, deleted_at
	`, username, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// MARK: Delete
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`, id)

	return err
}
