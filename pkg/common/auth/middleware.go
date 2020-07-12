package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

// Middleware extracts the auth token from the request, verifies it, and returns an updated context with the validated auth token.
// If the auth token is missing or invalid, the function writes a 403 Unauthorized status and error message to the response.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := parseAuthToken(r)
		if err != nil {
			errMsg := fmt.Sprintf("error parsing auth token from request: %s", err)
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		token, err := verifyAuthToken(ctx, authToken)
		if err != nil {
			errMsg := fmt.Sprintf("error verifying auth token: %s", err)
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		ctx = WithToken(ctx, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing auth token")
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 {
		return "", errors.New("invalid/malformed auth token")
	}
	authToken := splitAuthHeader[1]
	return authToken, nil
}

func verifyAuthToken(ctx context.Context, authToken string) (*Token, error) {
	// Get a reference to the Auth client
	client, err := getFireBaseApp().Auth(ctx)
	if err != nil {
		return nil, err
	}
	return verifyAuthTokenWithClient(client, ctx, authToken)
}

func verifyAuthTokenWithClient(client *auth.Client, ctx context.Context, authToken string) (*Token, error) {
	// Verify the auth token
	firebaseToken, err := client.VerifyIDTokenAndCheckRevoked(ctx, authToken)
	if err != nil {
		return nil, err
	}
	token := Token(*firebaseToken)
	return &token, nil
}
