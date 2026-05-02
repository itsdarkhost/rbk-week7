package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Username, req.Email, req.Password, "user")
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"access_token": token})
}
