package controller

import (
	"html/template"
	"net/http"
	"time"

	"vibecheck/internal/controller/view"
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
	entry, err := h.entrySvc.GetTodayEntry(userID, time.Now())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := h.tmpl.ExecuteTemplate(w, "home", view.HomeView{TodayEntry: entry}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
