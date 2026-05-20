package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"vibecheck/internal/domain"
)

type EntryRepo struct{ db *sql.DB }

func NewEntryRepo(db *sql.DB) *EntryRepo { return &EntryRepo{db: db} }

func (r *EntryRepo) CreateEntry(e *domain.Entry) (*domain.Entry, error) {
	out := &domain.Entry{}
	err := r.db.QueryRow(
		`INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, note)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, user_id, date, depression, happiness, pain, energy, sleep, note, created_at`,
		e.UserID, e.Date, e.Depression, e.Happiness, e.Pain, e.Energy, e.Sleep, e.Note,
	).Scan(&out.ID, &out.UserID, &out.Date, &out.Depression, &out.Happiness, &out.Pain, &out.Energy, &out.Sleep, &out.Note, &out.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, domain.ErrDuplicateEntry
		}
		return nil, fmt.Errorf("create entry: %w", err)
	}
	return out, nil
}
