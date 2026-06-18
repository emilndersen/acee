package albums

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	albums, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if albums == nil {
		albums = []Album{}
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"albums": albums,
	})
}

func (h *Handler) BySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteError(w, http.StatusBadRequest, "slug is required")
		return
	}

	album, err := h.repo.BySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "album not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"album": album,
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateAlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if input.Slug == "" {
		response.WriteError(w, http.StatusBadRequest, "slug is required")
		return
	}
	if input.Title == "" {
		response.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	album, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok":    true,
		"album": album,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteError(w, http.StatusBadRequest, "slug is required")
		return
	}

	var input UpdateAlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	album, err := h.repo.Update(r.Context(), slug, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "album not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"album": album,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteError(w, http.StatusBadRequest, "slug is required")
		return
	}

	err := h.repo.Delete(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "album not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
