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
			next.ServeHTTP(w, setContext(r, ""))
		}

		session, err := auth.CheckSession(cookie.Value)
		log.Println(session.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Errorf("cookie not loaded %s", err))

			next.ServeHTTP(w, setContext(r, ""))
		}

		if session.Expires.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Errorf("session expired"))
			next.ServeHTTP(w, setContext(r, ""))
		}

		next.ServeHTTP(w, setContext(r, session.ID))
	})
}

func setContext(r *http.Request, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), AuthCtx{}, value))
}
