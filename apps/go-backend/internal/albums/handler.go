package albums

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler хранит ссылку на repo.
type Handler struct {
	// repo нужен, чтобы ходить в базу.
	repo *Repo
}

// NewHandler создаёт новый albums handler.
func NewHandler(repo *Repo) *Handler {
	// Возвращаем handler с подключённым repo.
	return &Handler{repo: repo}
}

// List — GET /api/albums
// Возвращает список всех альбомов.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Получаем альбомы из базы.
	albums, err := h.repo.List(r.Context())
	if err != nil {
		// Если ошибка — отвечаем 500.
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Если база пустая, вернём пустой массив, а не null.
	if albums == nil {
		albums = []Album{}
	}

	// Отдаём JSON-ответ.
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":     true,
		"albums": albums,
	})
}

// BySlug — GET /api/albums/{slug}
// Возвращает один альбом по slug.
func (h *Handler) BySlug(w http.ResponseWriter, r *http.Request) {
	// Читаем slug из URL.
	slug := chi.URLParam(r, "slug")

	// Если slug пустой — это bad request.
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	// Ищем альбом в базе.
	album, err := h.repo.BySlug(r.Context(), slug)
	if err != nil {
		// Пока упрощённо отвечаем 500.
		// Позже можно сделать отдельную обработку not found.
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Отдаём найденный альбом.
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":    true,
		"album": album,
	})
}

// Create — POST /api/albums
// Создаёт новый альбом.
// Потом эту ручку надо будет защитить админской авторизацией.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Подготавливаем структуру под входной JSON.
	var input CreateAlbumInput

	// Декодируем JSON из тела запроса.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// Проверяем обязательные поля.
	if input.Slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	// Проверяем обязательные поля.
	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	// Создаём альбом в базе.
	album, err := h.repo.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Возвращаем созданный альбом.
	writeJSON(w, http.StatusCreated, map[string]any{
		"ok":    true,
		"album": album,
	})
}

// writeJSON — маленький helper для ответа JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	// Ставим правильный content-type.
	w.Header().Set("Content-Type", "application/json")

	// Ставим HTTP status code.
	w.WriteHeader(status)

	// Кодируем ответ в JSON.
	_ = json.NewEncoder(w).Encode(v)
}

// writeError — helper для ошибок в едином формате.
func writeError(w http.ResponseWriter, status int, msg string) {
	// Возвращаем JSON с ok=false и текстом ошибки.
	writeJSON(w, status, map[string]any{
		"ok":    false,
		"error": msg,
	})
}
