package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsdarkhost/rbk-week4/internal/services"
)

type Handler struct {
	userService    *services.UserService
	cityService    *services.CityService
	weatherService *services.WeatherService
}

// MARK: New Handler
func NewHandler(userService *services.UserService, cityService *services.CityService, weatherService *services.WeatherService) *Handler {
	return &Handler{
		userService:    userService,
		cityService:    cityService,
		weatherService: weatherService,
	}
}

// MARK: Routes
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.health)
	r.Post("/users", h.createUser)
	r.Get("/users", h.listUsers)
	r.Get("/users/{id}", h.getUser)
	r.Put("/users/{id}", h.updateUser)
	r.Delete("/users/{id}", h.deleteUser)
	r.Post("/users/{id}/cities", h.createCity)
	r.Get("/users/{id}/cities", h.listCities)
	r.Delete("/users/{id}/cities/{city_id}", h.deleteCity)
	r.Get("/users/{id}/weather", h.getWeather)
	r.Get("/users/{id}/weather/history", h.weatherHistory)

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
