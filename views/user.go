package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/tnyie/journaler-api/auth"
	"github.com/tnyie/journaler-api/middleware"
	"github.com/tnyie/journaler-api/models"
)

// GetCurrentUser
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{
		ID: r.Context().Value(middleware.AuthCtx{}).(string),
	}

	err := user.Get()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed to get user data %s", err))
		return
	}

	encoded, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed to marshal user object %s", err))
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}

// CreateUser creates a new internal user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("couldn't get json body"))
		return
	}

	var user *models.User
	var login *models.Login

	err = json.Unmarshal(bd, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("couldn't parse json body"))
		return
	}

	err = json.Unmarshal(bd, &login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// return json with offending issue
		log.Println(fmt.Errorf("failed to parse login data %s", err))
		return
	}

	// create user and check for errors
	if user.Create() != nil {
		w.WriteHeader(http.StatusConflict)
		// return json with conflicting fields
		log.Println(fmt.Errorf("error creating user %s", err))
		return
	}

	userAuth := &models.UserAuth{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Verified: false,
		Hash:     auth.HashPassword(login.Password),
	}

	// create internal user authentication model and check for errors
	if userAuth.Create() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed to create internal user auth %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Created new user ", user.Username)
}

// UpdateUser ...
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != r.Context().Value(middleware.AuthCtx{}).(string) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("user not authorized"))
		return
	}

	user := &models.User{
		ID: id,
	}

	err := user.Get()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("failed getting user %s", err))
		return
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("failed getting json %s", err))
		return
	}

	var newUser *models.User

	err = json.Unmarshal(bd, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("failed parsing json %s", err))
		return
	}
	newUser.ID = user.ID
}

// DeleteUser ...
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != r.Context().Value(middleware.AuthCtx{}).(string) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("user not authorized"))
		return
	}

	user := &models.User{
		ID: id,
	}

	user.Get()

	user.Delete()
}
