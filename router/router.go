package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tnyie/journaler-api/auth"
	"github.com/tnyie/journaler-api/middleware"
	"github.com/tnyie/journaler-api/views"
)

func Route(r *chi.Mux) {
	r.Use(middleware.AuthMiddleware)

	r.Mount("/auth", authHandler())

	r.Mount("/users", userHandler())
	r.Mount("/journals", journalHandler())
	r.Mount("/entries", entryHandler())

}

func authHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/{provider}", auth.RedirectURI)
	r.Get("/{provider}/callback", auth.Callback)

	return r
}

func userHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", views.GetCurrentUser)
	r.Post("/", views.CreateUser)
	r.Patch("/{id}", views.UpdateUser)
	r.Delete("/{id}", views.DeleteUser)

	return r
}

func journalHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", views.GetOwnJournals)
	r.Get("/{id}", views.GetJournalInfo)
	r.Post("/", views.CreateJournal)

	return r
}

func entryHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Get("/{id}", views.GetEntry)
	r.Post("/", views.PostEntry)

	return r
}
