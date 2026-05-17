package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	vibedb "vibecheck/db"

	"vibecheck/internal/controller"
	"vibecheck/internal/middleware"
	"vibecheck/internal/repo"
	"vibecheck/internal/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env:", err)
		}
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("open db:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("ping db:", err)
	}
	if err := vibedb.Migrate(db, "db/migrations"); err != nil {
		log.Fatal("migrate:", err)
	}

	tmpl := template.Must(template.New("").ParseGlob("web/templates/*.html"))
	template.Must(tmpl.ParseGlob("web/templates/pages/*.html"))
	template.Must(tmpl.ParseGlob("web/templates/components/*.html"))

	userRepo := repo.NewUserRepo(db)
	sessionRepo := repo.NewSessionRepo(db)
	authSvc := service.NewAuth(userRepo, sessionRepo)

	home := controller.NewHome(tmpl)
	auth := controller.NewAuth(tmpl, authSvc)

	requireAuth := middleware.RequireAuth(sessionRepo)

	mux := http.NewServeMux()
	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	mux.HandleFunc("GET /login", auth.LoginPage)
	mux.HandleFunc("POST /login", auth.Login)
	mux.HandleFunc("GET /register", auth.RegisterPage)
	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /logout", auth.Logout)

	mux.Handle("GET /", requireAuth(http.HandlerFunc(home.Index)))

	handler := middleware.Chain(mux)
	log.Println("listening on :8088")
	if err := http.ListenAndServe(":8088", handler); err != nil {
		log.Fatal(err)
	}
}
