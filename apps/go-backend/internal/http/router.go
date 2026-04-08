package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/emilndersen/acee/apps/go-backend/internal/albums"
	"github.com/emilndersen/acee/apps/go-backend/internal/bookings"
	"github.com/emilndersen/acee/apps/go-backend/internal/config"
	"github.com/emilndersen/acee/apps/go-backend/internal/photos"
	"github.com/emilndersen/acee/apps/go-backend/internal/users"
)

func NewRouter(pool *pgxpool.Pool, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://*.railway.app",
		},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	r.Get("/health/db", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "db": "connected"})
	})

	usersRepo := users.NewRepo(pool)
	usersHandler := users.NewHandler(usersRepo)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", usersHandler.List)
		r.Post("/", usersHandler.Create)
	})

	albumsRepo := albums.NewRepo(pool)
	albumsHandler := albums.NewHandler(albumsRepo)

	photosRepo := photos.NewRepo(pool)
	photosHandler := photos.NewHandler(photosRepo)

	r.Route("/api/albums", func(r chi.Router) {
		r.Get("/", albumsHandler.List)
		r.Get("/{slug}", albumsHandler.BySlug)
		r.Get("/{slug}/photos", photosHandler.ListByAlbumSlug)

		r.With(AdminOnly(cfg.AdminAPIToken)).Post("/", albumsHandler.Create)
		r.With(AdminOnly(cfg.AdminAPIToken)).Post("/{slug}/photos", photosHandler.CreateByAlbumSlug)
	})
	r.With(AdminOnly(cfg.AdminAPIToken)).Delete("/api/photos/{id}", photosHandler.Delete)

	bookingsRepo := bookings.NewRepo(pool)
	bookingsHandler := bookings.NewHandler(bookingsRepo)

	r.Route("/api/bookings", func(r chi.Router) {
		r.Post("/", bookingsHandler.Create)
		r.With(AdminOnly(cfg.AdminAPIToken)).Get("/", bookingsHandler.List)
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
