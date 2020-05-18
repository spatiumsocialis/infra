package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jinzhu/gorm"
)

type (
	// contextKey represents keys into a request context
	contextKey string
	// Token is a proxy of the Firebase Auth SDKs token, so that importing packages won't need to import both auth packages
	Token auth.Token

	// User is just a gorm model wrapper for the Firebase UID to support circle queries
	User struct {
		ID        string     `json:"id"`
		CircleID  string     `json:"circleId"`
		CreatedAt time.Time  `json:"-"`
		UpdatedAt time.Time  `json:"-"`
		DeletedAt *time.Time `json:"-"`
	}

	// Profile contains a user's profile information
	Profile struct {
		UID            string `json:"uid"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
	}
)

const (
	// TestUID is the UID of a test user
	TestUID = "UpIEj9XrQNMzdOQDgPSY0MGSsnO2"

	googleApplicationCredentialsKey = "GOOGLE_APPLICATION_CREDENTIALS"

	// contextKeyAuthToken is the key in context under which the auth token is stored in AuthenticateRequest
	contextKeyAuthToken = contextKey("auth-token")
)

var firebaseApp *firebase.App

// GetUser retrieves the UID from a request's token and returns the user model
func GetUser(r *http.Request, db *gorm.DB) (User, error) {
	// Retrieve the auth token from the request context
	token, err := GetTokenFrom(r.Context())
	if err != nil {
		return User{}, err
	}
	// Retrieve the client UID from the token
	currentUID := token.UID

	// Fetch the current user
	var user User
	if result := db.FirstOrCreate(&user, User{ID: currentUID}); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

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
		if err != nil {
			return []Profile{}, fmt.Errorf("error getting user '%v': %v", u.ID, err)
		}
		p.Name = userRecord.DisplayName
		p.ProfilePicture = userRecord.PhotoURL
		profiles[i] = p
	}
	return profiles, nil
}

// GetTokenFrom retrieves the access token from a request context. It returns an error if the token isn't found
func GetTokenFrom(ctx context.Context) (*Token, error) {
	token := ctx.Value(contextKeyAuthToken)
	if token == nil {
		err := errors.New("Error: context doesn't contain a token")
		return nil, err
	}
	t := Token(*(token.(*auth.Token)))
	return &t, nil
}

// AddTokenTo adds an access token to a request context
func AddTokenTo(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, contextKeyAuthToken, token)
}

// Middleware extracts the auth token from the request, verifies it, and returns an updated context with the validated auth token.
// If the auth token is missing or invalid, the function writes a 403 Unauthorized status and error message to the response.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the context from the request
		ctx := r.Context()
		// Retrieve the Authorization header from the request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := errors.New("Missing auth token")
			log.Printf("Error verifying token: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Split the Authorization header (since format should be "Bearer [TOKEN]")
		splitAuthHeader := strings.Split(authHeader, " ")
		if len(splitAuthHeader) != 2 {
			err := errors.New("Invalid/malformed auth token")
			log.Printf("Error verifying token: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Retrieve the authorization token
		authToken := splitAuthHeader[1]

		// Get a reference to the Auth client
		client, err := getFireBaseApp().Auth(ctx)
		if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
		}

		// Verify the auth token
		token, err := client.VerifyIDTokenAndCheckRevoked(ctx, authToken)
		if err != nil {
			log.Printf("Error verifying token: %v\n", err.Error())
			http.Error(w, "Error: Invalid auth token", http.StatusUnauthorized)
			return
		}

		// Add the token to the context and pass it to the next handler
		ctx = context.WithValue(ctx, contextKeyAuthToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
