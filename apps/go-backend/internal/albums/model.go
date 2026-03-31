package albums

// Album описывает одну запись альбома,
// которую мы будем возвращать на фронт.
type Album struct {
	// ID — UUID альбома из базы.
	ID string `json:"id"`

	// Slug — человекочитаемый идентификатор для URL.
	// Например: easytone-2026
	Slug string `json:"slug"`

	// Title — название альбома.
	Title string `json:"title"`

	// CoverURL — ссылка на обложку альбома.
	CoverURL string `json:"cover_url"`

	// Description — описание альбома.
	Description string `json:"description"`

	// IsPublic — публичен ли альбом.
	IsPublic bool `json:"is_public"`

	// SortOrder — порядок сортировки на главной.
	SortOrder int `json:"sort_order"`

	// CreatedAt — дата создания, приводим к строке для JSON.
	CreatedAt string `json:"created_at"`
}

// CreateAlbumInput описывает JSON,
// который мы ждём при создании альбома.
type CreateAlbumInput struct {
	// Slug обязателен, потому что по нему потом будем открывать страницу альбома.
	Slug string `json:"slug"`

	// Title — название альбома.
	Title string `json:"title"`

	// CoverURL — ссылка на картинку-обложку.
	CoverURL string `json:"cover_url"`

	// Description — короткое описание.
	Description string `json:"description"`

	// IsPublic — будет ли альбом виден всем.
	IsPublic bool `json:"is_public"`

	// SortOrder — порядок в сетке.
	SortOrder int `json:"sort_order"`
}
