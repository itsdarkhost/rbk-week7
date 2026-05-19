package main

import (
	"net/http"

	"github.com/itsdarkhost/rbk-week4/internal/clients"
	"github.com/itsdarkhost/rbk-week4/internal/config"
	"github.com/itsdarkhost/rbk-week4/internal/db"
	"github.com/itsdarkhost/rbk-week4/internal/handlers"
	"github.com/itsdarkhost/rbk-week4/internal/repos"
	"github.com/itsdarkhost/rbk-week4/internal/services"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	defer conn.Close()

	if err := db.Migrate(conn, cfg.MigrationsPath); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	userRepo := repos.NewUserRepo(conn)
	cityRepo := repos.NewCityRepo(conn)
	historyRepo := repos.NewWeatherHistoryRepo(conn)

	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	cityService := services.NewCityService(userRepo, cityRepo)
	weatherService := services.NewWeatherService(
		userRepo,
		cityRepo,
		historyRepo,
		clients.NewWeatherClient(cfg.WeatherAPIURL),
	)

	handler := handlers.NewHandler(userService, cityService, weatherService, cfg.JWTSecret, logger)

	logger.Info("server started", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, handler.Routes()); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
