package photos

type Photo struct {
	ID          string `json:"id"`
	AlbumID     string `json:"album_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ThumbURL    string `json:"thumb_url"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

type CreatePhotoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ThumbURL    string `json:"thumb_url"`
	SortOrder   int    `json:"sort_order"`
}
