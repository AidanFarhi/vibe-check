package main

import (
	"log"
	"net/http"
	"vibecheck/controller"
)

func main() {

	// instantiate controllers
	hc := controller.NewHomeController()

	// create multiplexer
	mux := http.NewServeMux()

	// static file server
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	// register handlers
	mux.HandleFunc("/", hc.GetHome)

	// create server
	server := http.Server{
		Addr:    ":8030",
		Handler: mux,
	}

	// start server
	log.Fatal(server.ListenAndServe())
}
