package albums

type Album struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	CategoryID   *string `json:"category_id"`
	CategoryName *string `json:"category_name,omitempty"`
	CoverURL     *string `json:"cover_url"`
	CoverPhotoID *string `json:"cover_photo_id"`
	IsPrivate    bool    `json:"is_private"`
	Password     *string `json:"-"`
	ViewCount    int     `json:"view_count"`
	CreatedAt    string  `json:"created_at"`
}

type CreateAlbumInput struct {
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	CategoryID  *string `json:"category_id"`
	IsPrivate   bool    `json:"is_private"`
	Password    *string `json:"password"`
}

type UpdateAlbumInput struct {
	Title        *string `json:"title"`
	Slug         *string `json:"slug"`
	Description  *string `json:"description"`
	CategoryID   *string `json:"category_id"`
	CoverPhotoID *string `json:"cover_photo_id"`
	IsPrivate    *bool   `json:"is_private"`
	Password     *string `json:"password"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	SortOrder int    `json:"sort_order"`
}

type AlbumAccessInput struct {
	Password string `json:"password"`
}
