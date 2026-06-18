package albums

type Album struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	CoverURL    string `json:"cover_url"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

type CreateAlbumInput struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	CoverURL    string `json:"cover_url"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateAlbumInput struct {
	Title       *string `json:"title"`
	CoverURL    *string `json:"cover_url"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"is_public"`
	SortOrder   *int    `json:"sort_order"`
}
