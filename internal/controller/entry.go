package controller

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

	_, err := c.entrySvc.SubmitEntry(
		userID,
		time.Now(),
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

	c.tmpl.ExecuteTemplate(w, "log-modal", view.ModalView{})
}
