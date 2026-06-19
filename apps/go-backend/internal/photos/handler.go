package photos

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListByAlbumSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	photos, err := h.repo.ListByAlbumSlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "photos": photos})
}

func (h *Handler) CreateByAlbumSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var input CreatePhotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.ImageURL == "" {
		response.WriteError(w, http.StatusBadRequest, "image_url required")
		return
	}
	photo, err := h.repo.CreateByAlbumSlug(r.Context(), slug, input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusCreated, map[string]any{"ok": true, "photo": photo})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input UpdatePhotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	photo, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "photo": photo})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusNotFound, "photo not found")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input struct {
		SortOrder int `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	order := input.SortOrder
	update := UpdatePhotoInput{SortOrder: &order}
	photo, err := h.repo.Update(r.Context(), id, update)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "photo": photo})
}
