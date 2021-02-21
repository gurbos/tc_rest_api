package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", APIHandler)
	r.HandleFunc("/productLines", ProductLineHandler).Name("productLines")
	r.HandleFunc("/{productLine:[a-zA-Z]+}/sets", CardSetHandler).Name("sets").Methods("GET")
	// r.HandleFunc("/{productLine:[a-zA-Z]+}/{set:[a-zA-Z-]+}/cards", CardHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
