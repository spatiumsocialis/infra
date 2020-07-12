package auth

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type (
	// contextKey represents keys into a request context
	contextKey string
	// Token is a proxy of the Firebase Auth SDKs token, so that importing packages won't need to import both auth packages explicitly
	Token auth.Token

	// Profile contains a user's profile information
	Profile struct {
		UID            string `json:"uid"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
	}
)

const (
	googleApplicationCredentialsKey = "GOOGLE_APPLICATION_CREDENTIALS"

	// contextKeyAuthToken is the key in context under which the auth token is stored in AuthenticateRequest
	contextKeyAuthToken = contextKey("auth-token")
)

var firebaseApp *firebase.App

// GetUserProfiles returns the users' profiles
func GetUserProfiles(users ...User) ([]Profile, error) {
	ctx := context.Background()
	client, err := getFireBaseApp().Auth(ctx)
	if err != nil {
		return []Profile{}, fmt.Errorf("error getting firebase app: %v", err)
	}
	profiles := make([]Profile, len(users))
	for i, u := range users {
		p := Profile{UID: u.ID}
		userRecord, err := client.GetUser(ctx, u.ID)
		if err == nil {
			p.Name = userRecord.DisplayName
			p.ProfilePicture = userRecord.PhotoURL
		}
		profiles[i] = p
	}
	return profiles, nil
}

func initializeApp() {

	if os.Getenv(googleApplicationCredentialsKey) == "" {
		log.Fatalf("auth: %v environment variable not set. Ensure you set %v to the path"+
			" to your service account JSON", googleApplicationCredentialsKey, googleApplicationCredentialsKey)
	}
	var err error
	firebaseApp, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}
}

func getFireBaseApp() *firebase.App {
	if firebaseApp == nil {
		initializeApp()
	}
	return firebaseApp
}
