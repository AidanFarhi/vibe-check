package controller

import (
	"html/template"
	"net/http"
)

type Home struct {
	tmpl *template.Template
}

func NewHome() *Home {
	tmpl := template.Must(template.New("").ParseGlob("web/templates/*.html"))
	template.Must(tmpl.ParseGlob("web/templates/pages/*.html"))
	template.Must(tmpl.ParseGlob("web/templates/components/*.html"))
	return &Home{tmpl: tmpl}
}

func (h *Home) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "home", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
