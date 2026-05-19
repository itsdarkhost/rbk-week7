package services

import (
	"context"
	"errors"
	"strings"

	"github.com/itsdarkhost/rbk-week4/internal/models"
)

var ErrEmptyCity = errors.New("city is required")

type CityRepository interface {
	Create(ctx context.Context, userId int, name string) (*models.City, error)
	List(ctx context.Context, userId int) ([]models.City, error)
	Delete(ctx context.Context, userId int, cityId int) error
}

type CityService struct {
	userRepo UserRepository
	cityRepo CityRepository
}

// MARK: New City Service
func NewCityService(userRepo UserRepository, cityRepo CityRepository) *CityService {
	return &CityService{userRepo: userRepo, cityRepo: cityRepo}
}

// MARK: Create
func (s *CityService) Create(ctx context.Context, userId int, name string) (*models.City, error) {
	_, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyCity
	}

	return s.cityRepo.Create(ctx, userId, name)
}

// MARK: List
func (s *CityService) List(ctx context.Context, userId int) ([]models.City, error) {
	_, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	return s.cityRepo.List(ctx, userId)
}

// MARK: Delete
func (s *CityService) Delete(ctx context.Context, userId int, cityId int) error {
	_, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}

	return s.cityRepo.Delete(ctx, userId, cityId)
}
