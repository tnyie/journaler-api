package views

import (
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/tnyie/journaler-api/middleware"
	"github.com/tnyie/journaler-api/models"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.AuthCtx{}).(*auth.Token)

	encoded, err := json.Marshal(&user)
	if err != nil {
		log.Println("Couldn't get user from token\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	respondJSON(w, encoded, http.StatusOK)
}

func EnableUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.AuthCtx{}).(*auth.Token).UID

	user := &models.User{
		ID: userID,
	}

	err := user.Get()

	if user.Enabled && err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = user.Enable()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't update user\n", err)
		return
	}

	journal := &models.Journal{
		Name:    "_default:" + userID,
		OwnerID: userID,
	}

	err = journal.Create()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Couldn't create user initial journal")
		return
	}
	w.WriteHeader(http.StatusCreated)
}
