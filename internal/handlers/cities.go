package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) createCity(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	var req struct {
		Name string `json:"name"`
		City string `json:"city"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	name := req.Name
	if name == "" {
		name = req.City
	}

	city, err := h.cityService.Create(r.Context(), user.Id, name)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, city)
}

func (h *Handler) listCities(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	cities, err := h.cityService.List(r.Context(), user.Id)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, cities)
}

func (h *Handler) deleteCity(w http.ResponseWriter, r *http.Request) {
	cityId, ok := readId(w, r, "city_id")
	if !ok {
		return
	}

	user, ok := UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	if err := h.cityService.Delete(r.Context(), user.Id, cityId); err != nil {
		writeRepoError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
