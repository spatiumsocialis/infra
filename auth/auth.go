package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jinzhu/gorm"
)

// TestUID is the UID of a test user
const TestUID = "UpIEj9XrQNMzdOQDgPSY0MGSsnO2"

const googleApplicationCredentialsKey = "GOOGLE_APPLICATION_CREDENTIALS"

type contextKey string

// ContextKeyAuthToken is the key in context under which the auth token is stored in AuthenticateRequest
const contextKeyAuthToken = contextKey("auth-token")

// Token is a proxy of the Firebase Auth SDKs token, so that importing packages won't need to import both auth packages
type Token auth.Token

var firebaseApp *firebase.App

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

// User is just a gorm model wrapper for the Firebase UID to support circle queries
type User struct {
	ID        string
	CircleID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

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
