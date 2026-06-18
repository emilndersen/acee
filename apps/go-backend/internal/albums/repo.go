package albums

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

const selectColumns = `id, slug, title, cover_url, description, is_public, sort_order, created_at::text`

func scanAlbum(row interface{ Scan(dest ...any) error }) (Album, error) {
	var a Album
	err := row.Scan(&a.ID, &a.Slug, &a.Title, &a.CoverURL, &a.Description, &a.IsPublic, &a.SortOrder, &a.CreatedAt)
	return a, err
}

func (r *Repo) List(ctx context.Context) ([]Album, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+selectColumns+`
		FROM albums
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Album
	for rows.Next() {
		a, err := scanAlbum(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) BySlug(ctx context.Context, slug string) (Album, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT `+selectColumns+`
		FROM albums WHERE slug = $1
	`, slug)
	return scanAlbum(row)
}

func (r *Repo) Create(ctx context.Context, in CreateAlbumInput) (Album, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO albums (slug, title, cover_url, description, is_public, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+selectColumns+`
	`, in.Slug, in.Title, in.CoverURL, in.Description, in.IsPublic, in.SortOrder)
	return scanAlbum(row)
}

func (r *Repo) Update(ctx context.Context, slug string, in UpdateAlbumInput) (Album, error) {
	var sets []string
	var args []any
	argN := 1

	if in.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", argN))
		args = append(args, *in.Title)
		argN++
	}
	if in.CoverURL != nil {
		sets = append(sets, fmt.Sprintf("cover_url = $%d", argN))
		args = append(args, *in.CoverURL)
		argN++
	}
	if in.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argN))
		args = append(args, *in.Description)
		argN++
	}
	if in.IsPublic != nil {
		sets = append(sets, fmt.Sprintf("is_public = $%d", argN))
		args = append(args, *in.IsPublic)
		argN++
	}
	if in.SortOrder != nil {
		sets = append(sets, fmt.Sprintf("sort_order = $%d", argN))
		args = append(args, *in.SortOrder)
		argN++
	}

	if len(sets) == 0 {
		return r.BySlug(ctx, slug)
	}

	args = append(args, slug)
	query := fmt.Sprintf(`
		UPDATE albums SET %s
		WHERE slug = $%d
		RETURNING %s
	`, strings.Join(sets, ", "), argN, selectColumns)

	row := r.pool.QueryRow(ctx, query, args...)
	return scanAlbum(row)
}

func (r *Repo) Delete(ctx context.Context, slug string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM albums WHERE slug = $1`, slug)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
