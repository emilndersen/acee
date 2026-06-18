package users

import (
	"encoding/json"
	"net/http"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"users": users,
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if input.Name == "" || input.Email == "" {
		response.WriteError(w, http.StatusBadRequest, "name and email required")
		return
	}

	user, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok":   true,
		"user": user,
	})
}
