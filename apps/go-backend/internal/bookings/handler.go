package bookings

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateBookingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Contact = strings.TrimSpace(input.Contact)
	input.Telegram = strings.TrimSpace(input.Telegram)
	input.ShootType = strings.TrimSpace(input.ShootType)
	input.Date = strings.TrimSpace(input.Date)
	input.Idea = strings.TrimSpace(input.Idea)

	if input.Name == "" {
		response.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if input.Contact == "" {
		response.WriteError(w, http.StatusBadRequest, "contact is required")
		return
	}
	if input.ShootType == "" {
		response.WriteError(w, http.StatusBadRequest, "shoot_type is required")
		return
	}

	booking, err := h.repo.Create(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		contact := booking.Contact
		if booking.Telegram != "" {
			contact = fmt.Sprintf("%s (TG: @%s)", booking.Contact, strings.TrimPrefix(booking.Telegram, "@"))
		}
		h.bot.SendBookingNotification(booking.Name, contact, booking.ShootType, booking.Date, booking.Idea)
	}()

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"ok":      true,
		"booking": booking,
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.repo.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if bookings == nil {
		bookings = []Booking{}
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"bookings": bookings,
	})
}

var validStatuses = map[string]bool{
	"new":       true,
	"confirmed": true,
	"completed": true,
	"cancelled": true,
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	var input UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Status = strings.TrimSpace(strings.ToLower(input.Status))
	if !validStatuses[input.Status] {
		response.WriteError(w, http.StatusBadRequest, "status must be one of: new, confirmed, completed, cancelled")
		return
	}

	booking, err := h.repo.UpdateStatus(r.Context(), id, input.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "booking not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"booking": booking,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "booking not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *Handler) BusyDates(w http.ResponseWriter, r *http.Request) {
	dates, err := h.repo.BusyDates(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if dates == nil {
		dates = []string{}
	}
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"dates": dates,
	})
}
