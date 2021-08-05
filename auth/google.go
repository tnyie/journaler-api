package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/tnyie/journaler-api/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig *oauth2.Config
)

func initGoogle() {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  viper.GetString("oauth.google.callback"),
		ClientID:     viper.GetString("oauth.google.id"),
		ClientSecret: viper.GetString("oauth.google.secret"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func beginGoogleAuth(w http.ResponseWriter, r *http.Request) {
	redirect_url := r.URL.Query().Get("redirect_url")

	state := uuid.New().String()
	rds.Set(context.Background(), state, redirect_url, time.Duration(time.Minute*3))

	url := googleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func completeGoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, code := r.FormValue("state"), r.FormValue("code")

	redirect_url := rds.Get(context.Background(), state)

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("code exchange failed %s", err))
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("failed getting user info %s", err))
		return
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user info response %s", err))
		return
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(contents, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user info response %s", err))
		return
	}

	if data["email"] == "" {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("couldn't get valid response from google %s", err))
		log.Println(data)
		return
	}

	user := &models.User{
		Email: data["email"].(string),
	}

	err = user.Get()

	if err != nil {
		user.Name = data["name"].(string)
		user.Username = data["name"].(string)
		user.External = true

		sanitary, field := user.Sanitize()
		if !sanitary {
			w.Write([]byte(field))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.Create()
	}

	cookie := CreateSession(user.ID)

	http.SetCookie(w, cookie)
	url, err := redirect_url.Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
