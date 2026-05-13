package controller

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"vibecheck/internal/domain"
	"vibecheck/internal/service"
)

type Auth struct {
	tmpl    *template.Template
	authSvc *service.AuthService
}

func NewAuth(tmpl *template.Template, authSvc *service.AuthService) *Auth {
	return &Auth{tmpl: tmpl, authSvc: authSvc}
}

type authView struct {
	Error string
}

func (a *Auth) RegisterPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, "register", authView{})
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	err := a.authSvc.Register(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		msg := "Registration failed. Please try again."
		switch {
		case errors.Is(err, service.ErrWeakPassword):
			msg = "Password must be at least 8 characters."
		case errors.Is(err, domain.ErrEmailTaken):
			msg = "An account with that email already exists."
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		a.render(w, "register", authView{Error: msg})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *Auth) LoginPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, "login", authView{})
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	token, err := a.authSvc.Login(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		a.render(w, "login", authView{Error: "Invalid email or password."})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		_ = a.authSvc.Logout(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *Auth) render(w http.ResponseWriter, name string, data authView) {
	if err := a.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
