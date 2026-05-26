package controller

import (
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
	v := view.HomeView{
		TodayEntry: entry,
		Chart:      buildChartView(entries),
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
			dep, hap, pain, energy, sleep := e.Depression, e.Happiness, e.Pain, e.Energy, e.Sleep
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
