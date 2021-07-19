package auth

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	Client *auth.Client
}

var AuthClient FirebaseClient

func InitFirebase() {
	ctx := context.Background()
	keyFilePath, err := filepath.Abs("/config/firebase.json")
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}
	opt := option.WithCredentialsFile(keyFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}
	AuthClient.Client, err = app.Auth(ctx)
	if err != nil {
		panic(fmt.Errorf("error initializing firebase authentication %v", err))
	}
}
