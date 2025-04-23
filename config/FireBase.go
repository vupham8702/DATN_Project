package config

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	AuthenClient *auth.Client
)

func InitializeFirebase() {
	opt := option.WithCredentialsFile(os.Getenv("URL_FIREBASE_JSON"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(fmt.Errorf("Initialize error : %v", err))
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Println(fmt.Errorf("Error getting Auth client: %v", err))
	}
	AuthenClient = authClient
}
