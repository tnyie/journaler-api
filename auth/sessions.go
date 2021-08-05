package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/tnyie/journaler-api/models"
)

var rds *redis.Client

var ctx = context.Background()

func CreateSession(userID string) *http.Cookie {
	id := uuid.New().String()

	expires := time.Now().Add(time.Hour * 24)
	session, err := json.Marshal(&models.UserSession{
		ID:      userID,
		Expires: expires,
	})
	if err != nil {
		return nil
	}

	duration := time.Until(expires)

	rds.Set(ctx, id, session, duration)

	return &http.Cookie{
		Name:     "sid",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
		// Secure:   true,
	}
}

func CheckSession(sessonID string) (*models.UserSession, error) {

	thing, err := rds.Get(ctx, sessonID).Result()
	if err != nil {
		return nil, err
	}

	var session *models.UserSession
	err = json.Unmarshal([]byte(thing), &session)

	return session, err
}
