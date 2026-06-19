package albums

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/emilndersen/acee/apps/go-backend/internal/photos"
	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

type Handler struct {
	repo      *Repo
	photoRepo *photos.Repo
	uploadDir string
}

func NewHandler(repo *Repo, photoRepo *photos.Repo, uploadDir string) *Handler {
	return &Handler{repo: repo, photoRepo: photoRepo, uploadDir: uploadDir}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var albumList []Album
	var err error

	if category != "" {
		albumList, err = h.repo.ListByCategory(r.Context(), category)
	} else {
		albumList, err = h.repo.List(r.Context())
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Hide private albums from public listing
	var visible []Album
	for _, a := range albumList {
		if !a.IsPrivate {
			visible = append(visible, a)
		}
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"albums": visible,
	})
}

func (h *Handler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	albumList, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"albums": albumList,
	})
}

func (h *Handler) BySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	album, err := h.repo.BySlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "album not found")
		return
	}

	if album.IsPrivate && album.Password != nil {
		pw := r.URL.Query().Get("password")
		if pw == "" {
			response.WriteJSON(w, http.StatusOK, map[string]any{
				"ok":             true,
				"requires_password": true,
				"album": map[string]any{
					"id":    album.ID,
					"title": album.Title,
					"slug":  album.Slug,
				},
			})
			return
		}
		if pw != *album.Password {
			response.WriteError(w, http.StatusForbidden, "wrong password")
			return
		}
	}

	go h.repo.IncrementViews(r.Context(), slug)

	album.Password = nil
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"album": album,
	})
}

func (h *Handler) CheckPassword(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var input AlbumAccessInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	album, err := h.repo.BySlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "album not found")
		return
	}

	if album.Password == nil || *album.Password != input.Password {
		response.WriteError(w, http.StatusForbidden, "wrong password")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateAlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Title == "" || input.Slug == "" {
		response.WriteError(w, http.StatusBadRequest, "title and slug required")
		return
	}
	album, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusCreated, map[string]any{"ok": true, "album": album})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var input UpdateAlbumInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	album, err := h.repo.Update(r.Context(), slug, input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "album": album})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if err := h.repo.Delete(r.Context(), slug); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.ListCategories(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "categories": cats})
}

func (h *Handler) DownloadZip(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	album, err := h.repo.BySlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "album not found")
		return
	}

	if album.IsPrivate && album.Password != nil {
		pw := r.URL.Query().Get("password")
		if pw != *album.Password {
			response.WriteError(w, http.StatusForbidden, "wrong password")
			return
		}
	}

	photoList, err := h.photoRepo.ListByAlbumSlug(r.Context(), slug)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, album.Slug))

	zw := zip.NewWriter(w)
	defer zw.Close()

	for _, p := range photoList {
		filePath := strings.TrimPrefix(p.ImageURL, "/uploads/")
		absPath := filepath.Join(h.uploadDir, filePath)

		f, err := os.Open(absPath)
		if err != nil {
			continue
		}

		ext := filepath.Ext(absPath)
		name := p.Title
		if name == "" {
			name = p.ID
		}
		name = name + ext

		entry, err := zw.Create(name)
		if err != nil {
			f.Close()
			continue
		}
		io.Copy(entry, f)
		f.Close()
	}
}
