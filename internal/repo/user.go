package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"vibecheck/internal/domain"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) CreateUser(email, passwordHash string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(
		`INSERT INTO "user" (email, password_hash) VALUES ($1, $2)
		 RETURNING id, email, password_hash, created_at`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, domain.ErrEmailTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, created_at FROM "user" WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
