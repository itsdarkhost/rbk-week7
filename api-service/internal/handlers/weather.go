package handlers

import (
	"errors"
	"net/http"

	"github.com/itsdarkhost/rbk-week4/internal/middleware"
)

func (h *Handler) getWeather(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	weather, err := h.weatherService.GetWeather(r.Context(), user.Id)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, weather)
}

func (h *Handler) weatherHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	limit := readQueryInt(r, "limit", 10)
	offset := readQueryInt(r, "offset", 0)
	city := r.URL.Query().Get("city")

	history, err := h.weatherService.History(r.Context(), user.Id, city, limit, offset)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, history)
}
