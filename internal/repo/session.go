package repo

import (
	"database/sql"
	"fmt"
	"time"

	"vibecheck/internal/domain"
)

type SessionRepo struct{ db *sql.DB }

func NewSessionRepo(db *sql.DB) *SessionRepo { return &SessionRepo{db: db} }

func (r *SessionRepo) CreateSession(userID, token string, expiresAt time.Time) (*domain.Session, error) {
	s := &domain.Session{}
	err := r.db.QueryRow(
		`INSERT INTO session (user_id, token, expires_at) VALUES ($1, $2, $3)
		 RETURNING id, user_id, token, expires_at, created_at`,
		userID, token, expiresAt,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return s, nil
}

func (r *SessionRepo) GetSessionByToken(token string) (*domain.Session, error) {
	s := &domain.Session{}
	err := r.db.QueryRow(
		`SELECT id, user_id, token, expires_at, created_at FROM session WHERE token = $1`,
		token,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SessionRepo) DeleteSession(token string) error {
	_, err := r.db.Exec(`DELETE FROM session WHERE token = $1`, token)
	return err
}
