package users

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

func (r *Repo) List(ctx context.Context) ([]User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, email, created_at::text
		FROM users
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}

	return out, rows.Err()
}

func (r *Repo) Create(ctx context.Context, in CreateUserInput) (User, error) {
	var u User

	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at::text
	`, in.Name, in.Email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
	)

	return u, err
}
