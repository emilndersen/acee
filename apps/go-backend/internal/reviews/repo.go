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

const selectColumns = `id, COALESCE(album_id::text, ''), author_name, text, rating, is_visible, created_at::text`

func (r *Repo) Create(ctx context.Context, in CreateReviewInput) (Review, error) {
	var rev Review
	var albumID *string
	if in.AlbumID != "" {
		albumID = &in.AlbumID
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO reviews (album_id, author_name, text, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING `+selectColumns+`
	`, albumID, in.AuthorName, in.Text, in.Rating).Scan(
		&rev.ID, &rev.AlbumID, &rev.AuthorName, &rev.Text, &rev.Rating, &rev.IsVisible, &rev.CreatedAt,
	)
	return rev, err
}

func (r *Repo) ListByAlbum(ctx context.Context, albumID string) ([]Review, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+selectColumns+`
		FROM reviews
		WHERE album_id = $1 AND is_visible = true
		ORDER BY created_at DESC
	`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Review
	for rows.Next() {
		var rev Review
		if err := rows.Scan(&rev.ID, &rev.AlbumID, &rev.AuthorName, &rev.Text, &rev.Rating, &rev.IsVisible, &rev.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rev)
	}
	return out, rows.Err()
}

func (r *Repo) ListVisible(ctx context.Context) ([]Review, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+selectColumns+`
		FROM reviews WHERE is_visible = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Review
	for rows.Next() {
		var rev Review
		if err := rows.Scan(&rev.ID, &rev.AlbumID, &rev.AuthorName, &rev.Text, &rev.Rating, &rev.IsVisible, &rev.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rev)
	}
	return out, rows.Err()
}

func (r *Repo) ListAll(ctx context.Context) ([]Review, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+selectColumns+`
		FROM reviews
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Review
	for rows.Next() {
		var rev Review
		if err := rows.Scan(&rev.ID, &rev.AlbumID, &rev.AuthorName, &rev.Text, &rev.Rating, &rev.IsVisible, &rev.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rev)
	}
	return out, rows.Err()
}

func (r *Repo) ToggleVisibility(ctx context.Context, id string) (Review, error) {
	var rev Review
	err := r.pool.QueryRow(ctx, `
		UPDATE reviews SET is_visible = NOT is_visible
		WHERE id = $1
		RETURNING `+selectColumns+`
	`, id).Scan(
		&rev.ID, &rev.AlbumID, &rev.AuthorName, &rev.Text, &rev.Rating, &rev.IsVisible, &rev.CreatedAt,
	)
	return rev, err
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
