package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tnyie/journaler-api/auth"
	"github.com/tnyie/journaler-api/models"
)

// CreateUser creates a new internal user
func Login(w http.ResponseWriter, r *http.Request) {
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("couldn't get json body"))
		return
	}

	var login *models.Login

	err = json.Unmarshal(bd, &login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("couldn't parse json body"))
		return
	}

	userAuth := &models.UserAuth{
		Email:    login.Email,
		Username: login.Username,
	}

	err = userAuth.Get()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("failed to get user information %s", err))
		return
	}

	if auth.CheckPassword(login.Password, userAuth.Hash) {
		http.SetCookie(w, auth.CreateSession(userAuth.ID))
		w.WriteHeader(http.StatusAccepted)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(fmt.Errorf("password did not match stored hash"))
		return
	}
}
