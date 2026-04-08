package photos

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

// ListByAlbumSlug возвращает все фото конкретного альбома по slug
func (r *Repo) ListByAlbumSlug(ctx context.Context, slug string) ([]Photo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			p.id,
			p.album_id,
			p.title,
			p.description,
			p.image_url,
			p.thumb_url,
			p.sort_order,
			p.created_at::text
		FROM photos_v2 p
		JOIN albums a ON a.id = p.album_id
		WHERE a.slug = $1
		ORDER BY p.sort_order ASC, p.created_at ASC
	`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Photo

	for rows.Next() {
		var p Photo

		if err := rows.Scan(
			&p.ID,
			&p.AlbumID,
			&p.Title,
			&p.Description,
			&p.ImageURL,
			&p.ThumbURL,
			&p.SortOrder,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	return out, rows.Err()
}

// CreateByAlbumSlug создаёт фото внутри альбома по slug
func (r *Repo) CreateByAlbumSlug(ctx context.Context, slug string, in CreatePhotoInput) (Photo, error) {
	var p Photo

	err := r.pool.QueryRow(ctx, `
		INSERT INTO photos_v2 (
			album_id,
			title,
			description,
			image_url,
			thumb_url,
			sort_order
		)
		VALUES (
			(SELECT id FROM albums WHERE slug = $1),
			$2, $3, $4, $5, $6
		)
		RETURNING
			id,
			album_id,
			title,
			description,
			image_url,
			thumb_url,
			sort_order,
			created_at::text
	`,
		slug,
		in.Title,
		in.Description,
		in.ImageURL,
		in.ThumbURL,
		in.SortOrder,
	).Scan(
		&p.ID,
		&p.AlbumID,
		&p.Title,
		&p.Description,
		&p.ImageURL,
		&p.ThumbURL,
		&p.SortOrder,
		&p.CreatedAt,
	)

	return p, err
}

// Delete удаляет фото по id
func (r *Repo) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `
		DELETE FROM photos_v2
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	// если ни одна строка не удалена, значит такого id не было
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
