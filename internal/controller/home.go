package controller

import (
	"html/template"
	"net/http"

	"vibecheck/internal/controller/view"
)

type Home struct {
	tmpl *template.Template
}

func NewHome(tmpl *template.Template) *Home {
	return &Home{tmpl: tmpl}
}

func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "home", view.ModalView{}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
