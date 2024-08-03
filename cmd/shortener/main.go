package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/handlers"
	"github.com/halviet/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.New()

	r := chi.NewRouter()

	r.Post("/", handlers.ShortenURLHandle(store))
	r.Get("/{id}", handlers.GetURLHandle(store))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
