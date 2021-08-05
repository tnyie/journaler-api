package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	"github.com/rs/cors"

	"github.com/tnyie/journaler-api/auth"
	"github.com/tnyie/journaler-api/config"
	"github.com/tnyie/journaler-api/database"
	"github.com/tnyie/journaler-api/models"
	"github.com/tnyie/journaler-api/router"
)

func main() {

	config.InitConfig()

	auth.InitAuth()

	goth.UseProviders(google.New(
		viper.GetString("providers.google.key"),
		viper.GetString("providers.google.secret"),
		"http://localhost:8080/auth/callback?provider=google",
		"email", "profile",
	))

	models.InitModels()
	database.InitDB()

	r := chi.NewRouter()

	router.Route(r)

	handler := cors.AllowAll().Handler(r)
	http.ListenAndServe(":8080", handler)
}
