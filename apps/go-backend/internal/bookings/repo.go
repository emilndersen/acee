package bookings

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

const selectColumns = `id, name, contact, telegram, shoot_type, date, shoot_date::text, idea, COALESCE(status, 'new'), created_at::text`

func (r *Repo) Create(ctx context.Context, in CreateBookingInput) (Booking, error) {
	var b Booking
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bookings (name, contact, telegram, shoot_type, date, shoot_date, idea)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+selectColumns+`
	`, in.Name, in.Contact, in.Telegram, in.ShootType, in.Date, in.ShootDate, in.Idea).Scan(
		&b.ID, &b.Name, &b.Contact, &b.Telegram, &b.ShootType, &b.Date, &b.ShootDate, &b.Idea, &b.Status, &b.CreatedAt,
	)
	return b, err
}

func (r *Repo) List(ctx context.Context) ([]Booking, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectColumns+` FROM bookings ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.Name, &b.Contact, &b.Telegram, &b.ShootType, &b.Date, &b.ShootDate, &b.Idea, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *Repo) BusyDates(ctx context.Context, month string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT shoot_date::text FROM bookings
		WHERE shoot_date IS NOT NULL
		AND to_char(shoot_date, 'YYYY-MM') = $1
		AND status != 'cancelled'
		ORDER BY 1
	`, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *Repo) UpdateStatus(ctx context.Context, id, status string) (Booking, error) {
	var b Booking
	err := r.pool.QueryRow(ctx, `
		UPDATE bookings SET status = $1 WHERE id = $2
		RETURNING `+selectColumns+`
	`, status, id).Scan(
		&b.ID, &b.Name, &b.Contact, &b.Telegram, &b.ShootType, &b.Date, &b.ShootDate, &b.Idea, &b.Status, &b.CreatedAt,
	)
	return b, err
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM bookings WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repo) BusyDates(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT date FROM bookings
		WHERE status IN ('new', 'confirmed')
		AND date != ''
		ORDER BY date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, rows.Err()
}
