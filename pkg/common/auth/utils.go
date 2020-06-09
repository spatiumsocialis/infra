package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// GenerateToken generates a valid Firebase auth token
func GenerateToken(uid string) (string, error) {

	// api key retrieved from https://console.firebase.google.com/u/0/project/spatiumsocialis-e4683/settings/general
	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		return "", errors.New("Error: GOOGLE_API_KEY env variable not set")
	}

	ctx := context.Background()
	client, err := getFireBaseApp().Auth(ctx)
	if err != nil {
		return "", err
	}
	// Generate a custom token based off the uid
	customToken, err := client.CustomToken(ctx, uid)
	if err != nil {
		return "", err
	}

	// Request an ID token using the custom token
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", apiKey)
	data := map[string]interface{}{"token": customToken, "returnSecureToken": true}
	payload, _ := json.Marshal(data)
	var res *http.Response
	if res, err = http.Post(url, "application/json", bytes.NewBuffer(payload)); err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New("Error requesting auth token: request failed")
	}

	// Retrieve the ID token from the response
	var bodyData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bodyData); err != nil {
		return "", err
	}
	token := bodyData["idToken"].(string)

	return token, nil
}
