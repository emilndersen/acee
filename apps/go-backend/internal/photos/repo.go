package photos

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

const selectColumns = `p.id, p.album_id, p.title, p.description, p.image_url, p.thumb_url, p.media_type, p.sort_order, p.created_at::text`

func scanPhoto(row interface{ Scan(dest ...any) error }) (Photo, error) {
	var p Photo
	err := row.Scan(&p.ID, &p.AlbumID, &p.Title, &p.Description, &p.ImageURL, &p.ThumbURL, &p.MediaType, &p.SortOrder, &p.CreatedAt)
	return p, err
}

func (r *Repo) ListByAlbumSlug(ctx context.Context, slug string) ([]Photo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.album_id, p.title, p.description, p.image_url, p.thumb_url, p.media_type, p.sort_order, p.created_at::text
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
		p, err := scanPhoto(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repo) CreateByAlbumSlug(ctx context.Context, slug string, in CreatePhotoInput) (Photo, error) {
	mediaType := in.MediaType
	if mediaType == "" {
		mediaType = "image"
	}

	row := r.pool.QueryRow(ctx, `
		INSERT INTO photos_v2 (album_id, title, description, image_url, thumb_url, media_type, sort_order)
		VALUES ((SELECT id FROM albums WHERE slug = $1), $2, $3, $4, $5, $6, $7)
		RETURNING `+selectColumns+`
	`, slug, in.Title, in.Description, in.ImageURL, in.ThumbURL, mediaType, in.SortOrder)
	return scanPhoto(row)
}

func (r *Repo) Update(ctx context.Context, id string, in UpdatePhotoInput) (Photo, error) {
	var sets []string
	var args []any
	argN := 1

	if in.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", argN))
		args = append(args, *in.Title)
		argN++
	}
	if in.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argN))
		args = append(args, *in.Description)
		argN++
	}
	if in.SortOrder != nil {
		sets = append(sets, fmt.Sprintf("sort_order = $%d", argN))
		args = append(args, *in.SortOrder)
		argN++
	}

	if len(sets) == 0 {
		row := r.pool.QueryRow(ctx, `SELECT `+selectColumns+` FROM photos_v2 WHERE id = $1`, id)
		return scanPhoto(row)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE photos_v2 SET %s
		WHERE id = $%d
		RETURNING %s
	`, strings.Join(sets, ", "), argN, selectColumns)

	row := r.pool.QueryRow(ctx, query, args...)
	return scanPhoto(row)
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM photos_v2 WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
