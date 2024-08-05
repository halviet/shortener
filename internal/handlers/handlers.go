package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/app"
	"github.com/halviet/shortener/internal/config"
	"github.com/halviet/shortener/internal/storage"
	"io"
	"net/http"
)

func ShortenURLHandle(store *storage.Store, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		resp, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if len(resp) < 11 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}

		urlID := app.RandString(8)
		url := cfg.BaseAddr + urlID

		store.SaveURL(storage.ShortURL{
			Origin: string(resp),
			Short:  urlID,
		})

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(url))
	}
}

func GetURLHandle(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "id")

		url, err := store.GetOrigin(urlID)
		if err != nil {
			http.NotFound(w, r)
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
