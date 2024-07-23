package main

import (
	"github.com/halviet/shortener/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", handlers.ShortenURLHandle)
	mux.HandleFunc("GET /{id}", handlers.GetURLHandle)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
