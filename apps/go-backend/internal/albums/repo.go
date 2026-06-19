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

const selectCols = `a.id, a.title, a.slug, a.description, a.category_id,
	c.name, COALESCE(p.thumb_url, p.image_url, ''),
	a.cover_photo_id, a.is_private, a.password, a.view_count, a.created_at::text`

const fromJoins = `FROM albums a
	LEFT JOIN album_categories c ON c.id = a.category_id
	LEFT JOIN photos_v2 p ON p.id = a.cover_photo_id`

func scanAlbum(row interface{ Scan(dest ...any) error }) (Album, error) {
	var a Album
	var coverURL *string
	err := row.Scan(&a.ID, &a.Title, &a.Slug, &a.Description, &a.CategoryID,
		&a.CategoryName, &coverURL,
		&a.CoverPhotoID, &a.IsPrivate, &a.Password, &a.ViewCount, &a.CreatedAt)
	a.CoverURL = coverURL
	return a, err
}

func (r *Repo) List(ctx context.Context) ([]Album, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectCols+` `+fromJoins+` ORDER BY a.created_at DESC`)
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
		a.Password = nil
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) ListByCategory(ctx context.Context, categorySlug string) ([]Album, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectCols+` `+fromJoins+`
		WHERE c.slug = $1 ORDER BY a.created_at DESC`, categorySlug)
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
		a.Password = nil
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) BySlug(ctx context.Context, slug string) (Album, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+selectCols+` `+fromJoins+` WHERE a.slug = $1`, slug)
	return scanAlbum(row)
}

func (r *Repo) IncrementViews(ctx context.Context, slug string) error {
	_, err := r.pool.Exec(ctx, `UPDATE albums SET view_count = view_count + 1 WHERE slug = $1`, slug)
	return err
}

func (r *Repo) Create(ctx context.Context, in CreateAlbumInput) (Album, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO albums (title, slug, description, category_id, is_private, password)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, slug, description, category_id,
			NULL::text, ''::text, cover_photo_id, is_private, password, view_count, created_at::text
	`, in.Title, in.Slug, in.Description, in.CategoryID, in.IsPrivate, in.Password)
	a, err := scanAlbum(row)
	if err == nil {
		a.Password = nil
	}
	return a, err
}

func (r *Repo) Update(ctx context.Context, slug string, in UpdateAlbumInput) (Album, error) {
	var sets []string
	var args []any
	n := 1

	if in.Title != nil {
		sets = append(sets, fmt.Sprintf("title = $%d", n))
		args = append(args, *in.Title)
		n++
	}
	if in.Slug != nil {
		sets = append(sets, fmt.Sprintf("slug = $%d", n))
		args = append(args, *in.Slug)
		n++
	}
	if in.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", n))
		args = append(args, *in.Description)
		n++
	}
	if in.CategoryID != nil {
		sets = append(sets, fmt.Sprintf("category_id = $%d", n))
		args = append(args, *in.CategoryID)
		n++
	}
	if in.CoverPhotoID != nil {
		sets = append(sets, fmt.Sprintf("cover_photo_id = $%d", n))
		args = append(args, *in.CoverPhotoID)
		n++
	}
	if in.IsPrivate != nil {
		sets = append(sets, fmt.Sprintf("is_private = $%d", n))
		args = append(args, *in.IsPrivate)
		n++
	}
	if in.Password != nil {
		sets = append(sets, fmt.Sprintf("password = $%d", n))
		args = append(args, *in.Password)
		n++
	}

	if len(sets) == 0 {
		return r.BySlug(ctx, slug)
	}

	args = append(args, slug)
	query := fmt.Sprintf(`UPDATE albums SET %s WHERE slug = $%d
		RETURNING id, title, slug, description, category_id,
			NULL::text, ''::text, cover_photo_id, is_private, password, view_count, created_at::text`,
		strings.Join(sets, ", "), n)

	a, err := scanAlbum(r.pool.QueryRow(ctx, query, args...))
	if err == nil {
		a.Password = nil
	}
	return a, err
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

func (r *Repo) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, slug, sort_order FROM album_categories ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
