package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type adminRow struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
	IsActive     bool
	CreatedAt    string
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (*adminRow, error) {
	var a adminRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, is_active, created_at::text
		FROM admins
		WHERE email = $1
	`, email).Scan(&a.ID, &a.Email, &a.PasswordHash, &a.DisplayName, &a.IsActive, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repo) CreateAdmin(ctx context.Context, email, password, displayName string) (Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Admin{}, err
	}

	var a Admin
	err = r.pool.QueryRow(ctx, `
		INSERT INTO admins (email, password_hash, display_name)
		VALUES ($1, $2, $3)
		RETURNING id, email, display_name, is_active, created_at::text
	`, email, string(hash), displayName).Scan(&a.ID, &a.Email, &a.DisplayName, &a.IsActive, &a.CreatedAt)
	return a, err
}

func (r *Repo) UpdatePassword(ctx context.Context, id, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE admins SET password_hash = $1 WHERE id = $2`, string(hash), id)
	return err
}
