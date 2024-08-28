package middleware

import (
	"github.com/halviet/shortener/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		rd *responseData
	}
)

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.rd.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.rd.status = statusCode
}

func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := loggingResponseWriter{
			ResponseWriter: w,
			rd: &responseData{
				status: 0,
				size:   0,
			},
		}

		next.ServeHTTP(&lw, r)

		logger.Log.Info(
			"Server response",
			zap.Int("status", lw.rd.status),
			zap.Int("size", lw.rd.size),
		)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Log.Info(
			"Incoming request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int64("duration", duration.Milliseconds()),
		)
	})
}
