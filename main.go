package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/tnyie/journaler-api/auth"
	"github.com/tnyie/journaler-api/config"
	"github.com/tnyie/journaler-api/database"
	"github.com/tnyie/journaler-api/models"
	"github.com/tnyie/journaler-api/router"
)

func main() {

	config.InitConfig()

	auth.InitFirebase()

	models.InitModels()

	database.InitDB()

	r := chi.NewRouter()

	router.Route(r)

	http.ListenAndServe(":8080", r)
}
