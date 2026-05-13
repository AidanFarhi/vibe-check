package controller

import (
	"html/template"
	"net/http"
)

type Home struct {
	tmpl *template.Template
}

func NewHome(tmpl *template.Template) *Home {
	return &Home{tmpl: tmpl}
}

func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "home", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
