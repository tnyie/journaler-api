package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/journaler-api/database"
)

func InitAuth() {

	rds = database.InitRedis()

	initGoogle()
	initGithub()
}

func RedirectURI(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	switch provider {
	case "google":
		beginGoogleAuth(w, r)
	case "github":
		beginGithubAuth(w, r)
	}
}

func Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	switch provider {
	case "google":
		completeGoogleAuth(w, r)
	case "github":
		completeGithubAuth(w, r)
	}
}
