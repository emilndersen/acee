package photos

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

// List возвращает все фото, отсортированные по sort_order
func (r *Repo) List(ctx context.Context) ([]Photo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, album, description, url, sort_order, created_at::text
		FROM photos
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Album,
			&p.Description, &p.URL,
			&p.SortOrder, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ByAlbum возвращает фото конкретного альбома
func (r *Repo) ByAlbum(ctx context.Context, album string) ([]Photo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, album, description, url, sort_order, created_at::text
		FROM photos
		WHERE album = $1
		ORDER BY sort_order ASC, created_at DESC
	`, album)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Album,
			&p.Description, &p.URL,
			&p.SortOrder, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Create добавляет новое фото (для админки)
func (r *Repo) Create(ctx context.Context, in CreatePhotoInput) (Photo, error) {
	var p Photo
	err := r.pool.QueryRow(ctx, `
		INSERT INTO photos (title, album, description, url, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, album, description, url, sort_order, created_at::text
	`, in.Title, in.Album, in.Description, in.URL, in.SortOrder).Scan(
		&p.ID, &p.Title, &p.Album,
		&p.Description, &p.URL,
		&p.SortOrder, &p.CreatedAt,
	)
	return p, err
}

// Delete удаляет фото по ID
func (r *Repo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM photos WHERE id = $1`, id)
	return err
}
