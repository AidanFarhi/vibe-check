package controller

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"vibecheck/internal/controller/view"
	"vibecheck/internal/domain"
	"vibecheck/internal/middleware"
	"vibecheck/internal/service"
)

type EntryController struct {
	tmpl     *template.Template
	entrySvc *service.EntryService
}

func NewEntry(tmpl *template.Template, entrySvc *service.EntryService) *EntryController {
	return &EntryController{tmpl: tmpl, entrySvc: entrySvc}
}

func (c *EntryController) Submit(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	parseInt := func(name string) int {
		v, _ := strconv.Atoi(r.FormValue(name))
		return v
	}

	entry, err := c.entrySvc.SubmitEntry(
		userID,
		localDateUTC(),
		parseInt("depression"),
		parseInt("happiness"),
		parseInt("pain"),
		parseInt("energy"),
		parseInt("sleep"),
		r.FormValue("note"),
	)
	if err != nil {
		msg := "Something went wrong. Please try again."
		if errors.Is(err, domain.ErrDuplicateEntry) {
			msg = "You've already logged an entry for today."
		}
		c.tmpl.ExecuteTemplate(w, "log-modal", view.ModalView{Open: true, Error: msg})
		return
	}

	entries, err := c.entrySvc.GetRecentEntries(userID, 7)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	streak, err := c.entrySvc.GetStreak(userID)
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
	c.tmpl.ExecuteTemplate(w, "entry-success", view.HomeView{
		TodayEntry: entry,
		Chart:      buildChartView(entries),
		Streak:     streak,
		Deltas:     buildMetricDeltas(entry, yesterdayEntry),
		ScoreLabel: buildScoreLabel(entry, yesterdayEntry),
	})
}
