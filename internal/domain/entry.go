package domain

import (
	"errors"
	"fmt"
	"time"
)

type Entry struct {
	ID         string
	UserID     string
	Date       time.Time
	Depression int
	Happiness  int
	Pain       int
	Energy     int
	Sleep      int
	Note       string
	CreatedAt  time.Time
	Score      float64
}

func (e *Entry) computeScore() float64 {
	return float64(e.Happiness+e.Energy+e.Sleep+(11-e.Depression)+(11-e.Pain)) / 5.0
}

func NewEntry(userID string, date time.Time, depression, happiness, pain, energy, sleep int, note string) (*Entry, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if date.IsZero() {
		return nil, errors.New("date is required")
	}
	for _, m := range []struct {
		name string
		val  int
	}{
		{"depression", depression},
		{"happiness", happiness},
		{"pain", pain},
		{"energy", energy},
		{"sleep", sleep},
	} {
		if m.val < 1 || m.val > 10 {
			return nil, fmt.Errorf("%s must be between 1 and 10", m.name)
		}
	}
	e := &Entry{
		UserID:     userID,
		Date:       date,
		Depression: depression,
		Happiness:  happiness,
		Pain:       pain,
		Energy:     energy,
		Sleep:      sleep,
		Note:       note,
	}
	e.Score = e.computeScore()
	return e, nil
}
