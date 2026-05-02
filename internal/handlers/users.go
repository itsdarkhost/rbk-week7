package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Create(r.Context(), req.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := readId(w, r, "id")
	if !ok {
		return
	}

	user, err := h.userService.Get(r.Context(), id)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := readId(w, r, "id")
	if !ok {
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Update(r.Context(), id, req.Username)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := readId(w, r, "id")
	if !ok {
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		writeRepoError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
