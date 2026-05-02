package handlers

import (
	"errors"
	"net/http"
)

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

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	current, ok := UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("user is required"))
		return
	}

	user, err := h.userService.Get(r.Context(), current.Id)
	if err != nil {
		writeRepoError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
