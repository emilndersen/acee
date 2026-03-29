package bookings

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Create — POST /api/bookings
// Принимает заявку на съёмку с фронта
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateBookingInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// валидация обязательных полей
	if input.Name == "" || input.Contact == "" {
		writeError(w, http.StatusBadRequest, "name and contact are required")
		return
	}

	booking, err := h.repo.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"ok":      true,
		"booking": booking,
	})
}

// List — GET /api/bookings (для админки)
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if bookings == nil {
		bookings = []Booking{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"bookings": bookings,
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
