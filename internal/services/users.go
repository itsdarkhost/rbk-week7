package services

import (
	"context"
	"strings"

	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/itsdarkhost/rbk-week4/internal/repos"
)

type UserService struct {
	repo *repos.UserRepo
}

// MARK: New User Service
func NewUserService(repo *repos.UserRepo) *UserService {
	return &UserService{repo: repo}
}

// MARK: List
func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

// MARK: Get
func (s *UserService) Get(ctx context.Context, id int) (*models.User, error) {
	return s.repo.Get(ctx, id)
}

// MARK: Create
func (s *UserService) Create(ctx context.Context, username string) (*models.User, error) {
	return s.repo.Create(ctx, strings.TrimSpace(username))
}

// MARK: Update
func (s *UserService) Update(ctx context.Context, id int, username string) (*models.User, error) {
	return s.repo.Update(ctx, id, strings.TrimSpace(username))
}

// MARK: Delete
func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
