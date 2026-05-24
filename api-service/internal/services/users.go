package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmptyEmail         = errors.New("email is required")
	ErrEmptyPassword      = errors.New("password is required")
	ErrEmailTaken         = errors.New("email already exists")
	ErrInvalidRole        = errors.New("invalid role")
)

type UserRepository interface {
	List(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, username string, email string, passwordHash string, role string) (*models.User, error)
	Update(ctx context.Context, id int, username string) (*models.User, error)
	Delete(ctx context.Context, id int) error
}

type UserService struct {
	repo      UserRepository
	jwtSecret []byte
}

// MARK: New User Service
func NewUserService(repo UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

// MARK: List
func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

// MARK: Get
func (s *UserService) Get(ctx context.Context, id int) (*models.User, error) {
	return s.repo.Get(ctx, id)
}

// MARK: Update
func (s *UserService) Update(ctx context.Context, id int, username string) (*models.User, error) {
	return s.repo.Update(ctx, id, strings.TrimSpace(username))
}

// MARK: Delete
func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// MARK: Register
func (s *UserService) Register(ctx context.Context, username string, email string, password string, role string) (*models.User, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	role = strings.TrimSpace(role)

	if email == "" {
		return nil, ErrEmptyEmail
	}
	if password == "" {
		return nil, ErrEmptyPassword
	}
	if username == "" {
		username = email
	}
	if role == "" {
		role = "user"
	}
	if role != "user" && role != "admin" {
		return nil, ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(ctx, username, email, string(hash), role)
	if isUniqueViolation(err) {
		return nil, ErrEmailTaken
	}

	return user, err
}

// MARK: Login
func (s *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
