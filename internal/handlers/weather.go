package handlers

import "net/http"

func (h *Handler) getWeather(w http.ResponseWriter, r *http.Request) {
	userId, ok := readId(w, r, "id")
	if !ok {
		return
	}

	weather, err := h.weatherService.GetWeather(r.Context(), userId)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, weather)
}

func (h *Handler) weatherHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := readId(w, r, "id")
	if !ok {
		return
	}

	limit := readQueryInt(r, "limit", 10)
	offset := readQueryInt(r, "offset", 0)
	city := r.URL.Query().Get("city")

	history, err := h.weatherService.History(r.Context(), userId, city, limit, offset)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, history)
}
