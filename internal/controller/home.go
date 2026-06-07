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
	estLoc, _ := time.LoadLocation("America/New_York")
	v := view.HomeView{
		TodayEntry:  entry,
		Chart:       buildChartView(entries, view.Period7D),
		Streak:      streak,
		Deltas:      buildMetricDeltas(entry, yesterdayEntry),
		ScoreLabel:  buildScoreLabel(entry, yesterdayEntry),
		CurrentDate: time.Now().In(estLoc).Format("Mon, Jan 2"),
	}
	if err := h.tmpl.ExecuteTemplate(w, "home", v); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Home) Chart(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	period := view.ParseChartPeriod(r.URL.Query().Get("period"))
	entries, err := h.entrySvc.GetRecentEntries(userID, view.DaysForPeriod(period))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	v := view.HomeView{Chart: buildChartView(entries, period)}
	if err := h.tmpl.ExecuteTemplate(w, "chart-card", v); err != nil {
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

func buildChartView(entries []domain.Entry, period view.ChartPeriod) view.ChartView {
	today := localDateUTC()
	byDate := make(map[string]domain.Entry, len(entries))
	for _, e := range entries {
		byDate[e.Date.Format("2006-01-02")] = e
	}
	var days []view.ChartDay
	switch period {
	case view.Period30D:
		days = dailyBuckets(today, 30, byDate)
	case view.Period3M:
		days = weeklyBuckets(today, 13, byDate)
	case view.Period6M:
		days = weeklyBuckets(today, 26, byDate)
	case view.Period1Y:
		days = monthlyBuckets(today, 12, byDate)
	default:
		days = dailyBuckets(today, 7, byDate)
	}
	return view.ChartView{Days: days, Period: period}
}

func dailyBuckets(today time.Time, count int, byDate map[string]domain.Entry) []view.ChartDay {
	days := make([]view.ChartDay, count)
	for i := 0; i < count; i++ {
		d := today.AddDate(0, 0, -(count - 1 - i))
		day := view.ChartDay{Label: dailyLabel(d, i, count)}
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
	return days
}

func dailyLabel(d time.Time, i, count int) string {
	if i == count-1 {
		return "Today"
	}
	if count == 7 {
		return d.Weekday().String()[:3]
	}
	// 30D: label every 5 days back from today.
	if (count-1-i)%5 == 0 {
		return d.Format("Jan 2")
	}
	return ""
}

func weeklyBuckets(today time.Time, count int, byDate map[string]domain.Entry) []view.ChartDay {
	days := make([]view.ChartDay, count)
	// Each bucket covers 7 days. Bucket count-1 ends on today.
	labelEvery := 2
	if count > 20 {
		labelEvery = 4
	}
	for i := 0; i < count; i++ {
		weeksAgo := count - 1 - i
		weekStart := today.AddDate(0, 0, -(weeksAgo*7 + 6))

		var sumScore float64
		var sumDep, sumHap, sumPain, sumEnergy, sumSleep int
		n := 0
		for offset := 0; offset < 7; offset++ {
			d := weekStart.AddDate(0, 0, offset)
			if e, ok := byDate[d.Format("2006-01-02")]; ok {
				sumScore += e.Score
				sumDep += e.Depression
				sumHap += e.Happiness
				sumPain += e.Pain
				sumEnergy += e.Energy
				sumSleep += e.Sleep
				n++
			}
		}

		day := view.ChartDay{Label: bucketLabel(weekStart.Format("Jan 2"), i, count, labelEvery)}
		if n > 0 {
			score := int(sumScore/float64(n) + 0.5)
			dep := roundDiv(sumDep, n)
			hap := roundDiv(sumHap, n)
			pain := roundDiv(sumPain, n)
			energy := roundDiv(sumEnergy, n)
			sleep := roundDiv(sumSleep, n)
			day.Score = &score
			day.Depression = &dep
			day.Happiness = &hap
			day.Pain = &pain
			day.Energy = &energy
			day.Sleep = &sleep
		}
		days[i] = day
	}
	return days
}

func monthlyBuckets(today time.Time, count int, byDate map[string]domain.Entry) []view.ChartDay {
	days := make([]view.ChartDay, count)
	for i := 0; i < count; i++ {
		monthsAgo := count - 1 - i
		target := today.AddDate(0, -monthsAgo, 0)
		bucketStart := time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, time.UTC)
		bucketEnd := bucketStart.AddDate(0, 1, 0)

		var sumScore float64
		var sumDep, sumHap, sumPain, sumEnergy, sumSleep int
		n := 0
		for d := bucketStart; d.Before(bucketEnd); d = d.AddDate(0, 0, 1) {
			if e, ok := byDate[d.Format("2006-01-02")]; ok {
				sumScore += e.Score
				sumDep += e.Depression
				sumHap += e.Happiness
				sumPain += e.Pain
				sumEnergy += e.Energy
				sumSleep += e.Sleep
				n++
			}
		}

		label := bucketStart.Format("Jan")
		if i == count-1 {
			label = "Today"
		}
		day := view.ChartDay{Label: label}
		if n > 0 {
			score := int(sumScore/float64(n) + 0.5)
			dep := roundDiv(sumDep, n)
			hap := roundDiv(sumHap, n)
			pain := roundDiv(sumPain, n)
			energy := roundDiv(sumEnergy, n)
			sleep := roundDiv(sumSleep, n)
			day.Score = &score
			day.Depression = &dep
			day.Happiness = &hap
			day.Pain = &pain
			day.Energy = &energy
			day.Sleep = &sleep
		}
		days[i] = day
	}
	return days
}

func bucketLabel(text string, i, count, every int) string {
	if i == count-1 {
		return "Today"
	}
	if (count-1-i)%every == 0 {
		return text
	}
	return ""
}

func roundDiv(sum, n int) int {
	return (sum*2 + n) / (2 * n)
}
