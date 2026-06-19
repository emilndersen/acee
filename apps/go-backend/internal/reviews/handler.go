package reviews

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
	"github.com/emilndersen/acee/apps/go-backend/internal/telegram"
)

type Handler struct {
	repo *Repo
	bot  *telegram.Bot
	pool *pgxpool.Pool
}

func NewHandler(repo *Repo, bot *telegram.Bot, pool *pgxpool.Pool) *Handler {
	return &Handler{repo: repo, bot: bot, pool: pool}
}

func (h *Handler) albumTitle(ctx context.Context, albumID string) string {
	if albumID == "" {
		return "—"
	}
	var title string
	err := h.pool.QueryRow(ctx, `SELECT title FROM albums WHERE id = $1`, albumID).Scan(&title)
	if err != nil {
		return "—"
	}
	return title
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.AuthorName = strings.TrimSpace(input.AuthorName)
	input.Text = strings.TrimSpace(input.Text)

	if input.AuthorName == "" {
		response.WriteError(w, http.StatusBadRequest, "author_name is required")
		return
	}
	if input.Text == "" {
		response.WriteError(w, http.StatusBadRequest, "text is required")
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		response.WriteError(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}

	review, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		title := h.albumTitle(context.Background(), review.AlbumID)
		h.bot.SendReviewNotification(review.AuthorName, review.Text, review.Rating, title)
	}()

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok":     true,
		"review": review,
	})
}

func (h *Handler) ListByAlbum(w http.ResponseWriter, r *http.Request) {
	albumID := chi.URLParam(r, "albumID")
	reviews, err := h.repo.ListByAlbum(r.Context(), albumID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if reviews == nil {
		reviews = []Review{}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"reviews": reviews,
	})
}

func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.repo.ListVisible(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if reviews == nil {
		reviews = []Review{}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"reviews": reviews,
	})
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.repo.ListAll(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if reviews == nil {
		reviews = []Review{}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"reviews": reviews,
	})
}

func (h *Handler) ToggleVisibility(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	review, err := h.repo.ToggleVisibility(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "review not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"review": review,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "review not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
