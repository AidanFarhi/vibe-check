package view

import "vibecheck/internal/domain"

type ModalView struct {
	Open  bool
	Error string
}

type ChartPeriod string

const (
	Period7D  ChartPeriod = "7D"
	Period30D ChartPeriod = "30D"
	Period3M  ChartPeriod = "3M"
	Period6M  ChartPeriod = "6M"
	Period1Y  ChartPeriod = "1Y"
)

func ParseChartPeriod(s string) ChartPeriod {
	switch ChartPeriod(s) {
	case Period30D, Period3M, Period6M, Period1Y:
		return ChartPeriod(s)
	}
	return Period7D
}

// DaysForPeriod returns how many days of entries to fetch to cover the period.
func DaysForPeriod(p ChartPeriod) int {
	switch p {
	case Period30D:
		return 30
	case Period3M:
		return 13 * 7
	case Period6M:
		return 26 * 7
	case Period1Y:
		return 366
	}
	return 7
}

// ParseChartMetric clamps the input to a known metric, defaulting to "score".
func ParseChartMetric(s string) string {
	switch s {
	case "energy", "sleep", "happiness", "pain", "depression", "score":
		return s
	}
	return "score"
}

// ChartDay holds data for a single bucket in the chart — a day, week, or month.
// Each metric pointer is nil when no entry exists for that bucket.
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
	Days   []ChartDay
	Period ChartPeriod
	Metric string
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
	TodayEntry  *domain.Entry
	Chart       ChartView
	Streak      int
	Deltas      MetricDeltas
	ScoreLabel  string
	CurrentDate string
}
