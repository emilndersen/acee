package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/emilndersen/acee/apps/go-backend/internal/bookings"
	"github.com/emilndersen/acee/apps/go-backend/internal/photos"
	"github.com/emilndersen/acee/apps/go-backend/internal/users"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// ── Middleware ──────────────────────────────────────────
	r.Use(middleware.Logger)    // логирует каждый запрос в консоль
	r.Use(middleware.Recoverer) // не падает если handler паникует

	// CORS — разрешает фронту на localhost:5173 обращаться к бэку
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // Vite dev server
			"https://*.railway.app", // продакшн
		},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	// ── Health ──────────────────────────────────────────────
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

	// ── Users ───────────────────────────────────────────────
	usersRepo := users.NewRepo(pool)
	usersHandler := users.NewHandler(usersRepo)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", usersHandler.List)
		r.Post("/", usersHandler.Create)
	})

	// ── Photos ──────────────────────────────────────────────
	photosRepo := photos.NewRepo(pool)
	photosHandler := photos.NewHandler(photosRepo)

	r.Route("/api/photos", func(r chi.Router) {
		r.Get("/", photosHandler.List)           // все фото
		r.Post("/", photosHandler.Create)        // добавить фото (админка)
		r.Get("/{album}", photosHandler.ByAlbum) // фото по альбому
		r.Delete("/{id}", photosHandler.Delete)  // удалить фото (админка)
	})

	// ── Bookings ────────────────────────────────────────────
	bookingsRepo := bookings.NewRepo(pool)
	bookingsHandler := bookings.NewHandler(bookingsRepo)

	r.Route("/api/bookings", func(r chi.Router) {
		r.Post("/", bookingsHandler.Create) // форма с фронта
		r.Get("/", bookingsHandler.List)    // список заявок (админка)
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
