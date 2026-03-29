package photos

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// List — GET /api/photos
// Возвращает все фото
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	photos, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// если фото ещё нет — вернём пустой массив, не null
	if photos == nil {
		photos = []Photo{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"photos": photos,
	})
}

// ByAlbum — GET /api/photos/{album}
// Возвращает фото конкретного альбома
func (h *Handler) ByAlbum(w http.ResponseWriter, r *http.Request) {
	album := chi.URLParam(r, "album") // читаем {album} из URL
	if album == "" {
		writeError(w, http.StatusBadRequest, "album is required")
		return
	}

	photos, err := h.repo.ByAlbum(r.Context(), album)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if photos == nil {
		photos = []Photo{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"album":  album,
		"photos": photos,
	})
}

// Create — POST /api/photos (только для админки)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreatePhotoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if input.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	photo, err := h.repo.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"ok":    true,
		"photo": photo,
	})
}

// Delete — DELETE /api/photos/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
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
