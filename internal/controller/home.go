package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"vibecheck/internal/controller/view"
	"vibecheck/internal/domain"
	"vibecheck/internal/middleware"
	"vibecheck/internal/service"
)

type Home struct {
	tmpl     *template.Template
	entrySvc *service.EntryService
}

func NewHome(tmpl *template.Template, entrySvc *service.EntryService) *Home {
	return &Home{tmpl: tmpl, entrySvc: entrySvc}
}

func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	entry, err := h.entrySvc.GetTodayEntry(userID, localDateUTC())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	entries, err := h.entrySvc.GetRecentEntries(userID, 7)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	streak, err := h.entrySvc.GetStreak(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	yesterday := localDateUTC().AddDate(0, 0, -1)
	var yesterdayEntry *domain.Entry
	for i := range entries {
		if entries[i].Date.Equal(yesterday) {
			yesterdayEntry = &entries[i]
			break
		}
	}
	v := view.HomeView{
		TodayEntry: entry,
		Chart:      buildChartView(entries),
		Streak:     streak,
		Deltas:     buildMetricDeltas(entry, yesterdayEntry),
		ScoreLabel: buildScoreLabel(entry, yesterdayEntry),
	}
	if err := h.tmpl.ExecuteTemplate(w, "home", v); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// localDateUTC returns the local calendar date as a UTC midnight time.Time.
// This ensures date comparisons match what PostgreSQL stores for DATE columns
// when the pq driver sends timestamps in UTC.
func localDateUTC() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func buildScoreLabel(today, yesterday *domain.Entry) string {
	if today == nil {
		return ""
	}
	tier := "Poor"
	switch {
	case today.Score >= 7.5:
		tier = "Great"
	case today.Score >= 4.5:
		tier = "Moderate"
	case today.Score >= 2.5:
		tier = "Low"
	}
	if yesterday == nil {
		return tier
	}
	switch {
	case today.Score > yesterday.Score:
		return tier + " ↑ up from yesterday"
	case today.Score < yesterday.Score:
		return tier + " ↓ down from yesterday"
	default:
		return tier + " same as yesterday"
	}
}

func buildMetricDeltas(today, yesterday *domain.Entry) view.MetricDeltas {
	if today == nil {
		return view.MetricDeltas{}
	}
	if yesterday == nil {
		noEntry := view.MetricDelta{Text: "— no entry yesterday"}
		return view.MetricDeltas{
			Energy: noEntry, Sleep: noEntry, Happiness: noEntry,
			Pain: noEntry, Depression: noEntry,
		}
	}
	return view.MetricDeltas{
		Energy:     metricDelta(today.Energy, yesterday.Energy, true),
		Sleep:      metricDelta(today.Sleep, yesterday.Sleep, true),
		Happiness:  metricDelta(today.Happiness, yesterday.Happiness, true),
		Pain:       metricDelta(today.Pain, yesterday.Pain, false),
		Depression: metricDelta(today.Depression, yesterday.Depression, false),
	}
}

func metricDelta(today, yesterday int, higherIsBetter bool) view.MetricDelta {
	d := today - yesterday
	if d == 0 {
		return view.MetricDelta{Text: "— no change from yesterday"}
	}
	abs, arrow := d, "↑"
	if d < 0 {
		abs, arrow = -d, "↓"
	}
	class := "good"
	if (d > 0) != higherIsBetter {
		class = "bad"
	}
	return view.MetricDelta{Text: fmt.Sprintf("%s %d from yesterday", arrow, abs), Class: class}
}

func buildChartView(entries []domain.Entry) view.ChartView {
	today := localDateUTC()
	byDate := make(map[string]domain.Entry, len(entries))
	for _, e := range entries {
		byDate[e.Date.Format("2006-01-02")] = e
	}
	days := make([]view.ChartDay, 7)
	for i := 0; i < 7; i++ {
		d := today.AddDate(0, 0, -(6 - i))
		label := d.Weekday().String()[:3]
		if i == 6 {
			label = "Today"
		}
		day := view.ChartDay{Label: label}
		if e, ok := byDate[d.Format("2006-01-02")]; ok {
			score := int(e.Score + 0.5)
			dep, hap, pain, energy, sleep := e.Depression, e.Happiness, e.Pain, e.Energy, e.Sleep
			day.Score = &score
			day.Depression = &dep
			day.Happiness = &hap
			day.Pain = &pain
			day.Energy = &energy
			day.Sleep = &sleep
		}
		days[i] = day
	}
	return view.ChartView{Days: days}
}
