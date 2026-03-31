package albums

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo хранит подключение к базе.
type Repo struct {
	// pool — общий пул соединений с Postgres.
	pool *pgxpool.Pool
}

// NewRepo создаёт новый albums repo.
func NewRepo(pool *pgxpool.Pool) *Repo {
	// Возвращаем готовую структуру с pool внутри.
	return &Repo{pool: pool}
}

// List возвращает все альбомы.
// Сортируем сначала по sort_order, потом по created_at.
func (r *Repo) List(ctx context.Context) ([]Album, error) {
	// Выполняем SQL-запрос.
	rows, err := r.pool.Query(ctx, `
		SELECT
			id,
			slug,
			title,
			cover_url,
			description,
			is_public,
			sort_order,
			created_at::text
		FROM albums
		ORDER BY sort_order ASC, created_at DESC
	`)
	if err != nil {
		// Если база вернула ошибку — отдаём её наверх.
		return nil, err
	}

	// После окончания чтения закрываем rows.
	defer rows.Close()

	// Здесь соберём все альбомы.
	var out []Album

	// Проходим по всем строкам результата.
	for rows.Next() {
		// Временная переменная под одну запись.
		var a Album

		// Сканируем поля из SQL в Go-структуру.
		if err := rows.Scan(
			&a.ID,
			&a.Slug,
			&a.Title,
			&a.CoverURL,
			&a.Description,
			&a.IsPublic,
			&a.SortOrder,
			&a.CreatedAt,
		); err != nil {
			// Если scan сломался — возвращаем ошибку.
			return nil, err
		}

		// Добавляем альбом в итоговый массив.
		out = append(out, a)
	}

	// Возвращаем массив и возможную ошибку итерации.
	return out, rows.Err()
}

// BySlug возвращает один альбом по его slug.
func (r *Repo) BySlug(ctx context.Context, slug string) (Album, error) {
	// Подготовим переменную под результат.
	var a Album

	// Ищем альбом по slug.
	err := r.pool.QueryRow(ctx, `
		SELECT
			id,
			slug,
			title,
			cover_url,
			description,
			is_public,
			sort_order,
			created_at::text
		FROM albums
		WHERE slug = $1
	`, slug).Scan(
		&a.ID,
		&a.Slug,
		&a.Title,
		&a.CoverURL,
		&a.Description,
		&a.IsPublic,
		&a.SortOrder,
		&a.CreatedAt,
	)

	// Возвращаем найденный альбом и ошибку, если была.
	return a, err
}

// Create создаёт новый альбом в базе.
func (r *Repo) Create(ctx context.Context, in CreateAlbumInput) (Album, error) {
	// Подготовим переменную под созданный альбом.
	var a Album

	// Выполняем INSERT и сразу возвращаем созданную запись.
	err := r.pool.QueryRow(ctx, `
		INSERT INTO albums (
			slug,
			title,
			cover_url,
			description,
			is_public,
			sort_order
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			slug,
			title,
			cover_url,
			description,
			is_public,
			sort_order,
			created_at::text
	`,
		in.Slug,
		in.Title,
		in.CoverURL,
		in.Description,
		in.IsPublic,
		in.SortOrder,
	).Scan(
		&a.ID,
		&a.Slug,
		&a.Title,
		&a.CoverURL,
		&a.Description,
		&a.IsPublic,
		&a.SortOrder,
		&a.CreatedAt,
	)

	// Возвращаем созданный альбом и ошибку.
	return a, err
}
