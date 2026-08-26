package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
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

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid auth format or token", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			hash := sha256.Sum256([]byte(token))
			hashStr := hex.EncodeToString(hash[:])

			record, err := store.GetTokenByHash(hashStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			log.Printf("AUTH DEBUG -> ID=%q Name=%q Namespace=%q Revoked=%t", record.ID, record.Name, record.Namespace, record.Revoked)

			if record.Revoked {
				http.Error(w, "token has been revoked", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), TokenContextKey, record)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
