package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsdarkhost/rbk-week4/internal/middleware"
	"github.com/itsdarkhost/rbk-week4/internal/services"
)

type Handler struct {
	userService    *services.UserService
	cityService    *services.CityService
	weatherService *services.WeatherService
	jwtSecret      []byte
}

// MARK: New Handler
func NewHandler(userService *services.UserService, cityService *services.CityService, weatherService *services.WeatherService, jwtSecret string) *Handler {
	return &Handler{
		userService:    userService,
		cityService:    cityService,
		weatherService: weatherService,
		jwtSecret:      []byte(jwtSecret),
	}
}

// MARK: Routes
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.health)
	r.Post("/auth/register", h.register)
	r.Post("/auth/login", h.login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtSecret))

		r.Get("/users/me", h.me)
		r.Post("/cities", h.createCity)
		r.Get("/cities", h.listCities)
		r.Delete("/cities/{city_id}", h.deleteCity)
		r.Get("/weather", h.getWeather)
		r.Get("/weather/history", h.weatherHistory)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/users", h.listUsers)
			r.Get("/users/{id}", h.getUser)
			r.Delete("/users/{id}", h.deleteUser)
		})
	})

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
