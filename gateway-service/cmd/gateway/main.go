package main

import (
	"net/http"

	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/clients"
	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/config"
	"github.com/itsdarkhost/rbk-week4/gateway-service/internal/handlers"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg := config.Load()
	handler := handlers.NewHandler(
		clients.NewWeatherClient(cfg.ExternalWeatherAPIURL),
		clients.NewAPIClient(cfg.APIServiceURL),
		logger,
	)

	logger.Info("gateway started", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, handler.Routes()); err != nil {
		logger.Fatal("gateway stopped", zap.Error(err))
	}
}
