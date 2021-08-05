package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tnyie/journaler-api/auth"
)

type AuthCtx struct{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sid")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Errorf("cookie not loaded %s", err))
			// login page redirect
			return
		}

		session, err := auth.CheckSession(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Errorf("cookie not loaded %s", err))
			// login page redirect
			return
		}

		if session.Expires.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Errorf("session expired"))
			// login page redirect
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AuthCtx{}, session.ID)))
	})
}
