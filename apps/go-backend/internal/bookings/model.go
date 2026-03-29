package bookings

type Booking struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Contact   string `json:"contact"`
	ShootType string `json:"shoot_type"`
	Date      string `json:"date"`
	Idea      string `json:"idea"`
	CreatedAt string `json:"created_at"`
}

type CreateBookingInput struct {
	Name      string `json:"name"`
	Contact   string `json:"contact"`
	ShootType string `json:"shoot_type"`
	Date      string `json:"date"`
	Idea      string `json:"idea"`
}
