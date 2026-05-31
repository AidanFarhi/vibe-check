package view

import "vibecheck/internal/domain"

type ModalView struct {
	Open  bool
	Error string
}

// ChartDay holds data for a single day slot in the 7-day chart.
// Each metric pointer is nil when no entry was logged for that day.
type ChartDay struct {
	Label      string
	Score      *int
	Depression *int
	Happiness  *int
	Pain       *int
	Energy     *int
	Sleep      *int
}

type ChartView struct {
	Days []ChartDay
}

type MetricDelta struct {
	Text  string
	Class string
}

type MetricDeltas struct {
	Energy, Sleep, Happiness, Pain, Depression MetricDelta
}

type HomeView struct {
	ModalView
	TodayEntry *domain.Entry
	Chart      ChartView
	Streak     int
	Deltas     MetricDeltas
	ScoreLabel string
}
