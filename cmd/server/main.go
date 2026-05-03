package main

import (
	"log"
	"net/http"

	"vibecheck/internal/controller"
	"vibecheck/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	home := controller.NewHome()
	mux.HandleFunc("GET /", home.Index)

	handler := middleware.Chain(mux)

	log.Println("listening on :8088")
	if err := http.ListenAndServe(":8088", handler); err != nil {
		log.Fatal(err)
	}
}
