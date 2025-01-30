package middleware

import (
	"net/http"
	"os"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	secret := os.Getenv("SERVER_API_KEY")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey != secret || apiKey == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
