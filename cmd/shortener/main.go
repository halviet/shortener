package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/config"
	"github.com/halviet/shortener/internal/handlers"
	"github.com/halviet/shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.New()

	store := storage.New()
	r := chi.NewRouter()

	r.Post("/", handlers.ShortenURLHandle(store, cfg))
	r.Get("/{id}", handlers.GetURLHandle(store))

	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Fatal(err)
	}
}
