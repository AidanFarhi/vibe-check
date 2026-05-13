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
	return &Entry{
		UserID:     userID,
		Date:       date,
		Depression: depression,
		Happiness:  happiness,
		Pain:       pain,
		Energy:     energy,
		Sleep:      sleep,
		Note:       note,
	}, nil
}
