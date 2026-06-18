package settings

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
	settings, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"settings": result,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var input UpdateSettingsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if len(input.Settings) == 0 {
		response.WriteError(w, http.StatusBadRequest, "settings required")
		return
	}

	if err := h.repo.Update(r.Context(), input.Settings); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	settings, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"settings": result,
	})
}
