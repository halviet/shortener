package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/app"
	"github.com/halviet/shortener/internal/config"
	"github.com/halviet/shortener/internal/logger"
	"github.com/halviet/shortener/internal/storage"
	"go.uber.org/zap"
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

type RequestCreateURL struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

func JSONShortenURLHandle(store *storage.Store, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req RequestCreateURL
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Log.Error("decoding request json", zap.Error(err))

			res, err := json.Marshal(ResponseErr{Error: "Bad Request"})
			if err != nil {
				logger.Log.Error("encoding error json response", zap.Error(err))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}

		if req.URL == "" {
			res, err := json.Marshal(ResponseErr{Error: "Bad Request"})
			if err != nil {
				logger.Log.Error("encoding error json response", zap.Error(err))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}

		urlID := app.RandString(8)
		url := cfg.BaseAddr + urlID

		store.SaveURL(storage.ShortURL{
			Origin: req.URL,
			Short:  urlID,
		})

		resp, err := json.Marshal(Response{Result: url})
		if err != nil {
			logger.Log.Error("encoding of response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
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
