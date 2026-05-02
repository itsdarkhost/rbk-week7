package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) createCity(w http.ResponseWriter, r *http.Request) {
	userId, ok := readId(w, r, "id")
	if !ok {
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

	city, err := h.cityService.Create(r.Context(), userId, name)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, city)
}

func (h *Handler) listCities(w http.ResponseWriter, r *http.Request) {
	userId, ok := readId(w, r, "id")
	if !ok {
		return
	}

	cities, err := h.cityService.List(r.Context(), userId)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, cities)
}

func (h *Handler) deleteCity(w http.ResponseWriter, r *http.Request) {
	userId, ok := readId(w, r, "id")
	if !ok {
		return
	}

	cityId, ok := readId(w, r, "city_id")
	if !ok {
		return
	}

	if err := h.cityService.Delete(r.Context(), userId, cityId); err != nil {
		writeRepoError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
