package settings

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

func (r *Repo) List(ctx context.Context) ([]Setting, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT key, value, updated_at::text FROM site_settings ORDER BY key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Setting
	for rows.Next() {
		var s Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repo) Update(ctx context.Context, settings map[string]string) error {
	for key, value := range settings {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO site_settings (key, value, updated_at)
			VALUES ($1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = now()
		`, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
