package main

import (
	"log"
	"net/http"

	"github.com/itsdarkhost/rbk-week4/internal/clients"
	"github.com/itsdarkhost/rbk-week4/internal/config"
	"github.com/itsdarkhost/rbk-week4/internal/db"
	"github.com/itsdarkhost/rbk-week4/internal/handlers"
	"github.com/itsdarkhost/rbk-week4/internal/repos"
	"github.com/itsdarkhost/rbk-week4/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.Migrate(conn, cfg.MigrationsPath); err != nil {
		log.Fatal(err)
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

	handler := handlers.NewHandler(userService, cityService, weatherService, cfg.JWTSecret)

	log.Printf("server started on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler.Routes()))
}
