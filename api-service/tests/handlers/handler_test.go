package handlers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsdarkhost/rbk-week4/internal/handlers"
	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/itsdarkhost/rbk-week4/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) List(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	users, _ := args.Get(0).([]models.User)
	return users, args.Error(1)
}

func (m *mockUserRepository) Get(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *mockUserRepository) Create(ctx context.Context, username string, email string, passwordHash string, role string) (*models.User, error) {
	args := m.Called(ctx, username, email, passwordHash, role)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, id int, username string) (*models.User, error) {
	args := m.Called(ctx, id, username)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *mockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestHandlerGetUser(t *testing.T) {
	const jwtSecret = "secret"

	t.Run("returns user JSON", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)
		repo.On("Get", mock.Anything, 7).
			Return(&models.User{Id: 7, Username: "admin", Email: "admin@example.com", Role: "admin"}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/7", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken(t, jwtSecret))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"id":7,"username":"admin","email":"admin@example.com","role":"admin"}`, rec.Body.String())
		repo.AssertExpectations(t)
	})

	t.Run("validates id", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)

		req := httptest.NewRequest(http.MethodGet, "/users/bad", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken(t, jwtSecret))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error":"invalid id"}`, rec.Body.String())
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
	})

	t.Run("maps missing user to not found", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)
		repo.On("Get", mock.Anything, 99).Return((*models.User)(nil), sql.ErrNoRows).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken(t, jwtSecret))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"error":"sql: no rows in result set"}`, rec.Body.String())
		repo.AssertExpectations(t)
	})
}

func TestHandlerRegister(t *testing.T) {
	const jwtSecret = "secret"

	t.Run("creates user", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)
		expected := &models.User{Id: 3, Username: "Alice", Email: "alice@example.com", Role: "user"}
		repo.On("Create",
			mock.Anything,
			"Alice",
			"alice@example.com",
			mock.MatchedBy(func(hash string) bool {
				return bcrypt.CompareHashAndPassword([]byte(hash), []byte("pass123")) == nil
			}),
			"user",
		).Return(expected, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"username":"Alice","email":"alice@example.com","password":"pass123"}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id":3,"username":"Alice","email":"alice@example.com","role":"user"}`, rec.Body.String())
		repo.AssertExpectations(t)
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"email":`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, readError(t, rec.Body.String()), "unexpected EOF")
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("validates email", func(t *testing.T) {
		repo := new(mockUserRepository)
		router := testRouter(repo, jwtSecret)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"username":"Alice","email":"","password":"pass123"}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"error":"email is required"}`, rec.Body.String())
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func testRouter(repo services.UserRepository, jwtSecret string) http.Handler {
	userService := services.NewUserService(repo, jwtSecret)
	return handlers.NewHandler(userService, nil, nil, jwtSecret, zap.NewNop()).Routes()
}

func adminToken(t *testing.T, secret string) string {
	t.Helper()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"email":   "admin@example.com",
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(secret))
	require.NoError(t, err)

	return token
}

func readError(t *testing.T, body string) string {
	t.Helper()

	var response map[string]string
	require.NoError(t, json.Unmarshal([]byte(body), &response))
	return response["error"]
}
