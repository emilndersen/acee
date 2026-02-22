package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/emilndersen/acee/apps/go-backend/internal/users"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// health
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	r.Get("/health/db", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"ok": true,
			"db": map[string]any{"ok": 1},
		})
	})

	// users
	repo := users.NewRepo(pool)
	handler := users.NewHandler(repo)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
	})

	return r
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
