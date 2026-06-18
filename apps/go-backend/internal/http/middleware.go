package httpapi

import (
	"net/http"
	"strings"

	"github.com/emilndersen/acee/apps/go-backend/internal/auth"
	"github.com/emilndersen/acee/apps/go-backend/internal/response"
)

func AdminOnly(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := auth.ValidateToken(jwtSecret, tokenStr)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := auth.ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
