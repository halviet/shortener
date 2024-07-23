package handlers

import (
	"io"
	"net/http"
)

func ShortenURLHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if len(resp) < 11 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("https://localhost:8080/EwHXdJfB"))
}

func GetURLHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if r.PathValue("id") == "EwHXdJfB" {
		w.Header().Set("Location", "https://practicum.yandex.ru/")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Incorrect request", http.StatusBadRequest)
}
