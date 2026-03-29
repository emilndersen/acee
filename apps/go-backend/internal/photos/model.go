package photos

type Photo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Album       string `json:"album"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

type CreatePhotoInput struct {
	Title       string `json:"title"`
	Album       string `json:"album"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sort_order"`
}
