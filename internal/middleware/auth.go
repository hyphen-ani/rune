package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"rune/internal/storage"
	"strings"
)

func AuthMiddleware(store storage.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Public Endpoint
			if r.URL.Path == "/unseal" || r.URL.Path == "/status" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				http.Error(w, "invalid auth format or token", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			hash := sha256.Sum256([]byte(token))
			hashStr := hex.EncodeToString(hash[:])

			_, err := store.GetTokenByHash(hashStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
