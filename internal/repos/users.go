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
		SELECT id, username, email, password_hash, role, deleted_at
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
		SELECT id, username, email, password_hash, role, deleted_at
		FROM users
		WHERE id = $1
		AND deleted_at IS NULL
	`, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// MARK: Get By Email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `
		SELECT id, username, email, password_hash, role, deleted_at
		FROM users
		WHERE email = $1
		AND deleted_at IS NULL
	`, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// MARK: Create
func (r *UserRepo) Create(ctx context.Context, username string, email string, passwordHash string, role string) (*models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, password_hash, role, deleted_at
	`, username, email, passwordHash, role)

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
		RETURNING id, username, email, password_hash, role, deleted_at
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
