package httpapi

import (
	"net/http"
	"time"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	convertLimiter := newRateLimiter(20, time.Minute)
	runLimiter := newRateLimiter(6, time.Minute)

	mux.Handle("POST /convert", rateLimitMiddleware(http.HandlerFunc(ConvertHandler), convertLimiter))
	mux.Handle("POST /run", rateLimitMiddleware(http.HandlerFunc(RunHandler), runLimiter))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
