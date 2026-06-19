package httpapi

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/emilndersen/acee/apps/go-backend/internal/albums"
	"github.com/emilndersen/acee/apps/go-backend/internal/auth"
	"github.com/emilndersen/acee/apps/go-backend/internal/bookings"
	"github.com/emilndersen/acee/apps/go-backend/internal/config"
	"github.com/emilndersen/acee/apps/go-backend/internal/photos"
	"github.com/emilndersen/acee/apps/go-backend/internal/response"
	"github.com/emilndersen/acee/apps/go-backend/internal/settings"
	"github.com/emilndersen/acee/apps/go-backend/internal/telegram"
	"github.com/emilndersen/acee/apps/go-backend/internal/upload"
)

func NewRouter(pool *pgxpool.Pool, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"https://*.railway.app",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	r.Get("/health/db", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "db": "connected"})
	})

	uploadDir := cfg.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	absUpload, _ := filepath.Abs(uploadDir)
	fs := http.StripPrefix("/uploads/", http.FileServer(http.Dir(absUpload)))
	r.Get("/uploads/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	// Auth
	authRepo := auth.NewRepo(pool)
	authHandler := auth.NewHandler(authRepo, cfg.JWTSecret)

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.With(AdminOnly(cfg.JWTSecret)).Get("/me", authHandler.Me)
	})

	// Upload
	uploadHandler := upload.NewHandler(absUpload)
	r.With(AdminOnly(cfg.JWTSecret)).Post("/api/upload", uploadHandler.Upload)

	// Albums
	albumsRepo := albums.NewRepo(pool)
	photosRepo := photos.NewRepo(pool)
	albumsHandler := albums.NewHandler(albumsRepo, photosRepo, absUpload)
	photosHandler := photos.NewHandler(photosRepo)

	r.Get("/api/categories", albumsHandler.ListCategories)

	r.Route("/api/albums", func(r chi.Router) {
		r.Get("/", albumsHandler.List)
		r.Get("/{slug}", albumsHandler.BySlug)
		r.Post("/{slug}/access", albumsHandler.CheckPassword)
		r.Get("/{slug}/photos", photosHandler.ListByAlbumSlug)
		r.Get("/{slug}/download", albumsHandler.DownloadZip)

		r.With(AdminOnly(cfg.JWTSecret)).Get("/admin/list", albumsHandler.ListAdmin)
		r.With(AdminOnly(cfg.JWTSecret)).Post("/", albumsHandler.Create)
		r.With(AdminOnly(cfg.JWTSecret)).Put("/{slug}", albumsHandler.Update)
		r.With(AdminOnly(cfg.JWTSecret)).Delete("/{slug}", albumsHandler.Delete)
		r.With(AdminOnly(cfg.JWTSecret)).Post("/{slug}/photos", photosHandler.CreateByAlbumSlug)
	})

	r.With(AdminOnly(cfg.JWTSecret)).Put("/api/photos/{id}", photosHandler.Update)
	r.With(AdminOnly(cfg.JWTSecret)).Delete("/api/photos/{id}", photosHandler.Delete)
	r.With(AdminOnly(cfg.JWTSecret)).Patch("/api/photos/{id}/order", photosHandler.UpdateOrder)

	// Telegram bot
	bot := telegram.NewBot(cfg.TelegramBotToken, cfg.TelegramChatID)

	// Bookings
	bookingsRepo := bookings.NewRepo(pool)
	bookingsHandler := bookings.NewHandler(bookingsRepo, bot)

	r.Route("/api/bookings", func(r chi.Router) {
		r.Post("/", bookingsHandler.Create)
		r.With(AdminOnly(cfg.JWTSecret)).Get("/", bookingsHandler.List)
		r.With(AdminOnly(cfg.JWTSecret)).Patch("/{id}/status", bookingsHandler.UpdateStatus)
		r.With(AdminOnly(cfg.JWTSecret)).Delete("/{id}", bookingsHandler.Delete)
	})

	// Settings
	settingsRepo := settings.NewRepo(pool)
	settingsHandler := settings.NewHandler(settingsRepo)

	r.Route("/api/settings", func(r chi.Router) {
		r.Get("/", settingsHandler.List)
		r.With(AdminOnly(cfg.JWTSecret)).Put("/", settingsHandler.Update)
	})

	return r
}
