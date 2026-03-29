package bookings

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Create сохраняет новую заявку на съёмку
func (r *Repo) Create(ctx context.Context, in CreateBookingInput) (Booking, error) {
	var b Booking
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bookings (name, contact, shoot_type, date, idea)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, contact, shoot_type, date, idea, created_at::text
	`, in.Name, in.Contact, in.ShootType, in.Date, in.Idea).Scan(
		&b.ID, &b.Name, &b.Contact,
		&b.ShootType, &b.Date, &b.Idea,
		&b.CreatedAt,
	)
	return b, err
}

// List возвращает все заявки (для админки)
func (r *Repo) List(ctx context.Context) ([]Booking, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, contact, shoot_type, date, idea, created_at::text
		FROM bookings
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(
			&b.ID, &b.Name, &b.Contact,
			&b.ShootType, &b.Date, &b.Idea,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
