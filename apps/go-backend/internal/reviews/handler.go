package reviews

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
	"github.com/emilndersen/acee/apps/go-backend/internal/telegram"
)

type Handler struct {
	repo *Repo
	bot  *telegram.Bot
}

func NewHandler(repo *Repo, bot *telegram.Bot) *Handler {
	return &Handler{repo: repo, bot: bot}
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
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "reviews": reviews})
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
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "reviews": reviews})
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
		response.WriteError(w, http.StatusBadRequest, "author_name required")
		return
	}
	if input.Text == "" {
		response.WriteError(w, http.StatusBadRequest, "text required")
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		input.Rating = 5
	}

	review, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		stars := strings.Repeat("⭐", review.Rating)
		text := fmt.Sprintf(
			"🚨 <b>Новый отзыв!</b>\n\n"+
				"<b>Автор:</b> %s\n"+
				"<b>Оценка:</b> %s\n"+
				"<b>Текст:</b> %s",
			review.AuthorName, stars, review.Text,
		)
		h.bot.Send(text)
	}()

	response.WriteJSON(w, http.StatusCreated, map[string]any{"ok": true, "review": review})
}

func (h *Handler) SetVisible(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input struct {
		Visible bool `json:"is_visible"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := h.repo.SetVisible(r.Context(), id, input.Visible); err != nil {
		response.WriteError(w, http.StatusNotFound, "review not found")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusNotFound, "review not found")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
