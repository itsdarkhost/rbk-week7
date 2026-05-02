package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/itsdarkhost/rbk-week4/internal/services"
)

func readId(w http.ResponseWriter, r *http.Request, key string) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, key))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid id"))
		return 0, false
	}

	return id, true
}

func readQueryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return number
}

func writeRepoError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, services.ErrEmptyCity):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
