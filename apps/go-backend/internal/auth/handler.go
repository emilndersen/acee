package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

type Handler struct {
	repo      *Repo
	jwtSecret string
}

func NewHandler(repo *Repo, jwtSecret string) *Handler {
	return &Handler{repo: repo, jwtSecret: jwtSecret}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		response.WriteError(w, http.StatusBadRequest, "email and password required")
		return
	}

	admin, err := h.repo.FindByEmail(r.Context(), input.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !admin.IsActive {
		response.WriteError(w, http.StatusForbidden, "account is deactivated")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	a := Admin{
		ID:          admin.ID,
		Email:       admin.Email,
		DisplayName: admin.DisplayName,
		IsActive:    admin.IsActive,
		CreatedAt:   admin.CreatedAt,
	}

	accessToken, err := GenerateAccessToken(h.jwtSecret, a)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, err := GenerateRefreshToken(h.jwtSecret, a)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"admin": a,
		"tokens": TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if input.RefreshToken == "" {
		response.WriteError(w, http.StatusBadRequest, "refresh_token required")
		return
	}

	claims, err := ValidateToken(h.jwtSecret, input.RefreshToken)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	admin, err := h.repo.FindByEmail(r.Context(), claims.Email)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "admin not found")
		return
	}

	if !admin.IsActive {
		response.WriteError(w, http.StatusForbidden, "account is deactivated")
		return
	}

	a := Admin{
		ID:          admin.ID,
		Email:       admin.Email,
		DisplayName: admin.DisplayName,
		IsActive:    admin.IsActive,
		CreatedAt:   admin.CreatedAt,
	}

	accessToken, err := GenerateAccessToken(h.jwtSecret, a)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, err := GenerateRefreshToken(h.jwtSecret, a)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok": true,
		"tokens": TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	admin, err := h.repo.FindByEmail(r.Context(), claims.Email)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "admin not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"ok": true,
		"admin": Admin{
			ID:          admin.ID,
			Email:       admin.Email,
			DisplayName: admin.DisplayName,
			IsActive:    admin.IsActive,
			CreatedAt:   admin.CreatedAt,
		},
	})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if input.OldPassword == "" || input.NewPassword == "" {
		response.WriteError(w, http.StatusBadRequest, "old_password and new_password required")
		return
	}

	if len(input.NewPassword) < 6 {
		response.WriteError(w, http.StatusBadRequest, "new password must be at least 6 characters")
		return
	}

	admin, err := h.repo.FindByEmail(r.Context(), claims.Email)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "admin not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.OldPassword)); err != nil {
		response.WriteError(w, http.StatusUnauthorized, "wrong current password")
		return
	}

	if err := h.repo.UpdatePassword(r.Context(), admin.ID, input.NewPassword); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}
