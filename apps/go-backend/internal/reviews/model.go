package reviews

type Review struct {
	ID         string `json:"id"`
	AuthorName string `json:"author_name"`
	Text       string `json:"text"`
	Rating     int    `json:"rating"`
	IsVisible  bool   `json:"is_visible"`
	CreatedAt  string `json:"created_at"`
}

type CreateReviewInput struct {
	AuthorName string `json:"author_name"`
	Text       string `json:"text"`
	Rating     int    `json:"rating"`
}
