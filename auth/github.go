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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/tnyie/journaler-api/models"
)

var (
	githubOAuthConfig *oauth2.Config
)

func initGithub() {
	githubOAuthConfig = &oauth2.Config{
		RedirectURL:  viper.GetString("oauth.github.callback"),
		ClientID:     viper.GetString("oauth.github.id"),
		ClientSecret: viper.GetString("oauth.github.secret"),
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
}

func beginGithubAuth(w http.ResponseWriter, r *http.Request) {

	redirect_url := r.URL.Query().Get("redirect_url")

	state := uuid.New().String()
	rds.Set(ctx, state, redirect_url, time.Duration(time.Minute*3))

	url := githubOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type githubEmail struct {
	Email      string      `json:"email,omitempty"`
	Primary    bool        `json:"primary,omitempty"`
	Verified   bool        `json:"verified,omitempty"`
	Visibility interface{} `json:"visibility,omitempty"`
}

func completeGithubAuth(w http.ResponseWriter, r *http.Request) {
	state, code := r.FormValue("state"), r.FormValue("code")

	redirect_url := rds.Get(context.Background(), state)

	token, err := githubOAuthConfig.Exchange(ctx, code)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("code exchange failed %s", err))
		return
	}

	emailReq, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed to get user info %s", err))
		return
	}

	emailReq.Header.Set("Authorization", "token "+token.AccessToken)

	emailResponse, err := http.DefaultClient.Do(emailReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("failed getting user info %s", err))
		return
	}

	defer emailResponse.Body.Close()

	emailContents, err := ioutil.ReadAll(emailResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user info response %s", err))
		return
	}

	var emails []githubEmail

	err = json.Unmarshal(emailContents, &emails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user email response %s", err))
		return
	}

	userReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed to get user info %s", err))
		return
	}

	log.Println(token.AccessToken)

	userReq.Header.Set("Authorization", "token "+token.AccessToken)

	userResponse, err := http.DefaultClient.Do(userReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(fmt.Errorf("failed getting user info %s", err))
		return
	}

	defer userResponse.Body.Close()

	userContents, err := ioutil.ReadAll(userResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user info response %s", err))
		return
	}

	githubUser := make(map[string]interface{})
	err = json.Unmarshal(userContents, &githubUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("failed reading user info response %s", err))
		return
	}

	log.Println(githubUser)

	if len(emails) > 0 && githubUser["login"].(string) != "" && githubUser["name"] != "" {
		user := &models.User{
			Email: emails[0].Email,
		}

		err = user.Get()
		if err != nil {
			user.Name = githubUser["name"].(string)
			user.Username = githubUser["login"].(string)
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
}
