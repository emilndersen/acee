package reviews

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

const selectCols = `id, author_name, text, rating, is_visible, created_at::text`

func (r *Repo) ListVisible(ctx context.Context) ([]Review, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectCols+` FROM reviews WHERE is_visible = true ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *Repo) ListAll(ctx context.Context) ([]Review, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectCols+` FROM reviews ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *Repo) Create(ctx context.Context, in CreateReviewInput) (Review, error) {
	var rv Review
	err := r.pool.QueryRow(ctx, `
		INSERT INTO reviews (author_name, text, rating)
		VALUES ($1, $2, $3)
		RETURNING `+selectCols+`
	`, in.AuthorName, in.Text, in.Rating).Scan(
		&rv.ID, &rv.AuthorName, &rv.Text, &rv.Rating, &rv.IsVisible, &rv.CreatedAt,
	)
	return rv, err
}

func (r *Repo) SetVisible(ctx context.Context, id string, visible bool) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE reviews SET is_visible = $1 WHERE id = $2`, visible, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func scanRows(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]Review, error) {
	var out []Review
	for rows.Next() {
		var rv Review
		if err := rows.Scan(&rv.ID, &rv.AuthorName, &rv.Text, &rv.Rating, &rv.IsVisible, &rv.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rv)
	}
	return out, rows.Err()
}
