package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"vibecheck/internal/domain"
)

type EntryRepo struct{ db *sql.DB }

func NewEntryRepo(db *sql.DB) *EntryRepo { return &EntryRepo{db: db} }

func (r *EntryRepo) CreateEntry(e *domain.Entry) (*domain.Entry, error) {
	out := &domain.Entry{}
	err := r.db.QueryRow(
		`INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, note, score)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, user_id, date, depression, happiness, pain, energy, sleep, COALESCE(note, ''), score, created_at`,
		e.UserID, e.Date, e.Depression, e.Happiness, e.Pain, e.Energy, e.Sleep, e.Note, e.Score,
	).Scan(&out.ID, &out.UserID, &out.Date, &out.Depression, &out.Happiness, &out.Pain, &out.Energy, &out.Sleep, &out.Note, &out.Score, &out.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, domain.ErrDuplicateEntry
		}
		return nil, fmt.Errorf("create entry: %w", err)
	}
	return out, nil
}

func (r *EntryRepo) GetEntriesInRange(userID string, from, to time.Time) ([]domain.Entry, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, date, depression, happiness, pain, energy, sleep, COALESCE(note, ''), score, created_at
		 FROM entry WHERE user_id = $1 AND date BETWEEN $2 AND $3 ORDER BY date ASC`,
		userID, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("get entries in range: %w", err)
	}
	defer rows.Close()
	var entries []domain.Entry
	for rows.Next() {
		var e domain.Entry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Date, &e.Depression, &e.Happiness, &e.Pain, &e.Energy, &e.Sleep, &e.Note, &e.Score, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("get entries in range: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *EntryRepo) GetStreak(userID string, today time.Time) (int, error) {
	rows, err := r.db.Query(
		`SELECT date FROM entry WHERE user_id = $1 ORDER BY date DESC`,
		userID,
	)
	if err != nil {
		return 0, fmt.Errorf("get streak: %w", err)
	}
	defer rows.Close()
	var dates []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return 0, fmt.Errorf("get streak: %w", err)
		}
		dates = append(dates, d)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("get streak: %w", err)
	}
	if len(dates) == 0 {
		return 0, nil
	}
	yesterday := today.AddDate(0, 0, -1)
	start := dates[0]
	if !start.Equal(today) && !start.Equal(yesterday) {
		return 0, nil
	}
	streak := 0
	expected := start
	for _, d := range dates {
		if d.Equal(expected) {
			streak++
			expected = expected.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak, nil
}

func (r *EntryRepo) GetTodayEntry(userID string, date time.Time) (*domain.Entry, error) {
	out := &domain.Entry{}
	err := r.db.QueryRow(
		`SELECT id, user_id, date, depression, happiness, pain, energy, sleep, COALESCE(note, ''), score, created_at
		 FROM entry WHERE user_id = $1 AND date = $2`,
		userID, date,
	).Scan(&out.ID, &out.UserID, &out.Date, &out.Depression, &out.Happiness, &out.Pain, &out.Energy, &out.Sleep, &out.Note, &out.Score, &out.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get today entry: %w", err)
	}
	return out, nil
}
