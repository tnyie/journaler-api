package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/tnyie/journaler-api/auth"
)

type AuthCtx struct{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if parts[0] != "Bearer" || len(parts) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := auth.AuthClient.Client.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			log.Println("Invalid token\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthCtx{}, token)))
	})
}
