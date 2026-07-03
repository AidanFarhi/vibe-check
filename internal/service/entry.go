package service

import (
	"time"

	"vibecheck/internal/domain"
)

type EntryService struct {
	entries domain.EntryRepository
}

func NewEntry(entries domain.EntryRepository) *EntryService {
	return &EntryService{entries: entries}
}

func (s *EntryService) GetTodayEntry(userID string, date time.Time) (*domain.Entry, error) {
	return s.entries.GetTodayEntry(userID, date)
}

func (s *EntryService) GetRecentEntries(userID string, today time.Time, days int) ([]domain.Entry, error) {
	from := today.AddDate(0, 0, -(days - 1))
	return s.entries.GetEntriesInRange(userID, from, today)
}

func (s *EntryService) GetStreak(userID string, today time.Time) (int, error) {
	return s.entries.GetStreak(userID, today)
}

func (s *EntryService) SubmitEntry(userID string, date time.Time, depression, happiness, pain, energy, sleep int, note string) (*domain.Entry, error) {
	e, err := domain.NewEntry(userID, date, depression, happiness, pain, energy, sleep, note)
	if err != nil {
		return nil, err
	}
	return s.entries.CreateEntry(e)
}
