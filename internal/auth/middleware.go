package auth

import (
	"net/http"
	"os"
)

func ApiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		validAPIKey := os.Getenv("GOLD_API_KEY")
		if validAPIKey == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success":false,"error":"Server Configuration Error"}`))
			return
		}

		if r.Header.Get("X-API-Key") != validAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"ไม่มีสิทธิ์เข้าถึง (Invalid API Key)"}`))
			return
		}

		next(w, r)
	}
}