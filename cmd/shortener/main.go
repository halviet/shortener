package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/config"
	"github.com/halviet/shortener/internal/handlers"
	"github.com/halviet/shortener/internal/logger"
	mw "github.com/halviet/shortener/internal/middleware"
	"github.com/halviet/shortener/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err = logger.Init(cfg.LogLevel); err != nil {
		panic(err)
	}

	store := storage.New()
	r := chi.NewRouter()

	r.Use(mw.ResponseLogger)
	r.Use(mw.RequestLogger)
	r.Use(mw.GzipCompress)

	r.Post("/", handlers.ShortenURLHandle(store, cfg))
	r.Get("/{id}", handlers.GetURLHandle(store))

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.JSONShortenURLHandle(store, cfg))
	})

	logger.Log.Info("Running on", zap.String("addr", cfg.Addr))
	if err = http.ListenAndServe(cfg.Addr, r); err != nil {
		logger.Log.Fatal("Server crash", zap.String("error", err.Error()))
	}
}
