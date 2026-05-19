package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/itsdarkhost/rbk-week4/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

type mockCityRepository struct {
	mock.Mock
}

func (m *mockCityRepository) Create(ctx context.Context, userId int, name string) (*models.City, error) {
	args := m.Called(ctx, userId, name)
	city, _ := args.Get(0).(*models.City)
	return city, args.Error(1)
}

func (m *mockCityRepository) List(ctx context.Context, userId int) ([]models.City, error) {
	args := m.Called(ctx, userId)
	cities, _ := args.Get(0).([]models.City)
	return cities, args.Error(1)
}

func (m *mockCityRepository) Delete(ctx context.Context, userId int, cityId int) error {
	args := m.Called(ctx, userId, cityId)
	return args.Error(0)
}

type mockWeatherHistoryRepository struct {
	mock.Mock
}

func (m *mockWeatherHistoryRepository) Create(ctx context.Context, userId int, weather models.Weather) (*models.WeatherHistory, error) {
	args := m.Called(ctx, userId, weather)
	history, _ := args.Get(0).(*models.WeatherHistory)
	return history, args.Error(1)
}

func (m *mockWeatherHistoryRepository) List(ctx context.Context, userId int, city string, limit int, offset int) ([]models.WeatherHistory, error) {
	args := m.Called(ctx, userId, city, limit, offset)
	history, _ := args.Get(0).([]models.WeatherHistory)
	return history, args.Error(1)
}

type mockWeatherClient struct {
	mock.Mock
}

func (m *mockWeatherClient) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	args := m.Called(ctx, city)
	weather, _ := args.Get(0).(*models.Weather)
	return weather, args.Error(1)
}

func TestUserServiceRegister(t *testing.T) {
	t.Run("creates user with normalized input", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")
		expected := &models.User{Id: 1, Username: "Alice", Email: "alice@example.com", Role: "user"}

		repo.On("Create",
			mock.Anything,
			"Alice",
			"alice@example.com",
			mock.MatchedBy(func(hash string) bool {
				return bcrypt.CompareHashAndPassword([]byte(hash), []byte("pass123")) == nil
			}),
			"user",
		).Return(expected, nil).Once()

		user, err := service.Register(context.Background(), " Alice ", " ALICE@EXAMPLE.COM ", " pass123 ", "")

		require.NoError(t, err)
		assert.Equal(t, expected, user)
		repo.AssertExpectations(t)
	})

	t.Run("validates empty email", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")

		user, err := service.Register(context.Background(), "Alice", " ", "pass123", "user")

		require.ErrorIs(t, err, ErrEmptyEmail)
		assert.Nil(t, user)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("validates empty password", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")

		user, err := service.Register(context.Background(), "Alice", "alice@example.com", " ", "user")

		require.ErrorIs(t, err, ErrEmptyPassword)
		assert.Nil(t, user)
	})

	t.Run("validates invalid role", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")

		user, err := service.Register(context.Background(), "Alice", "alice@example.com", "pass123", "manager")

		require.ErrorIs(t, err, ErrInvalidRole)
		assert.Nil(t, user)
	})

	t.Run("returns repository error", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")
		repoErr := errors.New("insert failed")

		repo.On("Create", mock.Anything, "alice@example.com", "alice@example.com", mock.AnythingOfType("string"), "user").
			Return((*models.User)(nil), repoErr).Once()

		user, err := service.Register(context.Background(), "", "alice@example.com", "pass123", "")

		require.ErrorIs(t, err, repoErr)
		assert.Nil(t, user)
		repo.AssertExpectations(t)
	})
}

func TestUserServiceLogin(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	t.Run("returns signed token", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")
		repo.On("GetByEmail", mock.Anything, "alice@example.com").
			Return(&models.User{Id: 1, Email: "alice@example.com", PasswordHash: string(hash), Role: "admin"}, nil).Once()

		token, err := service.Login(context.Background(), " ALICE@EXAMPLE.COM ", "pass123")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		repo.AssertExpectations(t)
	})

	t.Run("rejects blank credentials", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")

		token, err := service.Login(context.Background(), "", "pass123")

		require.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		repo.AssertNotCalled(t, "GetByEmail", mock.Anything, mock.Anything)
	})

	t.Run("hides missing user", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")
		repo.On("GetByEmail", mock.Anything, "missing@example.com").
			Return((*models.User)(nil), sql.ErrNoRows).Once()

		token, err := service.Login(context.Background(), "missing@example.com", "pass123")

		require.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		repo.AssertExpectations(t)
	})

	t.Run("rejects wrong password", func(t *testing.T) {
		repo := new(mockUserRepository)
		service := NewUserService(repo, "secret")
		repo.On("GetByEmail", mock.Anything, "alice@example.com").
			Return(&models.User{Id: 1, Email: "alice@example.com", PasswordHash: string(hash), Role: "user"}, nil).Once()

		token, err := service.Login(context.Background(), "alice@example.com", "wrong")

		require.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		repo.AssertExpectations(t)
	})
}

func TestUserServiceSimpleMethods(t *testing.T) {
	repo := new(mockUserRepository)
	service := NewUserService(repo, "secret")
	expected := &models.User{Id: 5, Username: "updated"}

	repo.On("List", mock.Anything).Return([]models.User{{Id: 1}}, nil).Once()
	repo.On("Get", mock.Anything, 5).Return(&models.User{Id: 5}, nil).Once()
	repo.On("Update", mock.Anything, 5, "updated").Return(expected, nil).Once()
	repo.On("Delete", mock.Anything, 5).Return(nil).Once()

	users, err := service.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 1)

	user, err := service.Get(context.Background(), 5)
	require.NoError(t, err)
	assert.Equal(t, 5, user.Id)

	user, err = service.Update(context.Background(), 5, " updated ")
	require.NoError(t, err)
	assert.Equal(t, expected, user)

	require.NoError(t, service.Delete(context.Background(), 5))
	repo.AssertExpectations(t)
}

func TestCityService(t *testing.T) {
	t.Run("creates trimmed city after user exists", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		cityRepo := new(mockCityRepository)
		service := NewCityService(userRepo, cityRepo)

		userRepo.On("Get", mock.Anything, 1).Return(&models.User{Id: 1}, nil).Once()
		cityRepo.On("Create", mock.Anything, 1, "Almaty").
			Return(&models.City{Id: 10, UserId: 1, Name: "Almaty"}, nil).Once()

		city, err := service.Create(context.Background(), 1, " Almaty ")

		require.NoError(t, err)
		assert.Equal(t, "Almaty", city.Name)
		userRepo.AssertExpectations(t)
		cityRepo.AssertExpectations(t)
	})

	t.Run("rejects empty city", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		cityRepo := new(mockCityRepository)
		service := NewCityService(userRepo, cityRepo)

		userRepo.On("Get", mock.Anything, 1).Return(&models.User{Id: 1}, nil).Once()

		city, err := service.Create(context.Background(), 1, " ")

		require.ErrorIs(t, err, ErrEmptyCity)
		assert.Nil(t, city)
		cityRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestWeatherService(t *testing.T) {
	t.Run("gets weather and stores history", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		cityRepo := new(mockCityRepository)
		historyRepo := new(mockWeatherHistoryRepository)
		client := new(mockWeatherClient)
		service := NewWeatherService(userRepo, cityRepo, historyRepo, client)
		weather := models.Weather{City: "Almaty", Temperature: 21.5, Description: "clear"}
		history := &models.WeatherHistory{Id: 1, UserId: 1, City: "Almaty", Temperature: 21.5, Description: "clear", RequestedAt: time.Now()}

		userRepo.On("Get", mock.Anything, 1).Return(&models.User{Id: 1}, nil).Once()
		cityRepo.On("List", mock.Anything, 1).Return([]models.City{{Id: 10, UserId: 1, Name: "Almaty"}}, nil).Once()
		client.On("GetWeather", mock.Anything, "Almaty").Return(&weather, nil).Once()
		historyRepo.On("Create", mock.Anything, 1, weather).Return(history, nil).Once()

		items, err := service.GetWeather(context.Background(), 1)

		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, "Almaty", items[0].City)
		userRepo.AssertExpectations(t)
		cityRepo.AssertExpectations(t)
		client.AssertExpectations(t)
		historyRepo.AssertExpectations(t)
	})

	t.Run("normalizes history query parameters", func(t *testing.T) {
		userRepo := new(mockUserRepository)
		cityRepo := new(mockCityRepository)
		historyRepo := new(mockWeatherHistoryRepository)
		client := new(mockWeatherClient)
		service := NewWeatherService(userRepo, cityRepo, historyRepo, client)
		expected := []models.WeatherHistory{{Id: 1, City: "Almaty"}}

		userRepo.On("Get", mock.Anything, 1).Return(&models.User{Id: 1}, nil).Once()
		historyRepo.On("List", mock.Anything, 1, "Almaty", 10, 0).Return(expected, nil).Once()

		response, err := service.History(context.Background(), 1, " Almaty ", 0, -5)

		require.NoError(t, err)
		assert.Equal(t, expected, response.History)
		assert.Equal(t, "Almaty", response.City)
		userRepo.AssertExpectations(t)
		historyRepo.AssertExpectations(t)
	})
}
