package main

import (
	"log"
	"net/http"
	"os"

	"github.com/itsdarkhost/rbk-week4/internal/clients"
	"github.com/itsdarkhost/rbk-week4/internal/db"
	"github.com/itsdarkhost/rbk-week4/internal/handlers"
	"github.com/itsdarkhost/rbk-week4/internal/repos"
	"github.com/itsdarkhost/rbk-week4/internal/services"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		log.Fatal(err)
	}

	userRepo := repos.NewUserRepo(conn)
	cityRepo := repos.NewCityRepo(conn)
	historyRepo := repos.NewWeatherHistoryRepo(conn)

	userService := services.NewUserService(userRepo)
	cityService := services.NewCityService(userRepo, cityRepo)
	weatherService := services.NewWeatherService(
		userRepo,
		cityRepo,
		historyRepo,
		clients.NewWeatherClient(os.Getenv("WEATHER_API_URL")),
	)

	handler := handlers.NewHandler(userService, cityService, weatherService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler.Routes()))
}
