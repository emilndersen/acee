package photos

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// ListByAlbumSlug — GET /api/albums/{slug}/photos
func (h *Handler) ListByAlbumSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	photos, err := h.repo.ListByAlbumSlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if photos == nil {
		photos = []Photo{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"slug":   slug,
		"photos": photos,
	})
}

// CreateByAlbumSlug — POST /api/albums/{slug}/photos
func (h *Handler) CreateByAlbumSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	var input CreatePhotoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.ImageURL = strings.TrimSpace(input.ImageURL)
	input.ThumbURL = strings.TrimSpace(input.ThumbURL)

	if input.ImageURL == "" {
		writeError(w, http.StatusBadRequest, "image_url is required")
		return
	}

	photo, err := h.repo.CreateByAlbumSlug(r.Context(), slug, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"ok":    true,
		"slug":  slug,
		"photo": photo,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{
		"ok":    false,
		"error": msg,
	})
}

// Delete — DELETE /api/photos/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "photo not found")
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}
