package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/safe-distance/initialize"
)

func TestMain(m *testing.M) {
	initialize.Env(true, "test.env")
	os.Exit(m.Run())
}

func addTokenToRequest(r *http.Request, token string) {
	r.Header.Add("Authorization", "Bearer "+token)
}

func testMiddlewareHelper(t *testing.T, shouldSucceed bool, testToken func(validToken string) string) {
	// uid for matt@axial.technology user
	const uid = "UpIEj9XrQNMzdOQDgPSY0MGSsnO2"
	// api key retrieved from https://console.firebase.google.com/u/0/project/safe-distance-e4683/settings/general
	const apiKey = "AIzaSyBGtOPFjnB0nzjI_LoLmDrea-96xDckde4"

	ctx := context.Background()
	client, err := getFireBaseApp().Auth(ctx)
	if err != nil {
		t.Fatalf("Error getting Auth client: %v\n", err)
	}

	// Generate a custom token based off the uid
	customToken, err := client.CustomToken(ctx, uid)
	if err != nil {
		t.Fatalf("Error minting custom token: %v\n", err)
	}

	// Request an ID token using the custom token
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", apiKey)
	data := map[string]interface{}{"token": customToken, "returnSecureToken": true}
	payload, _ := json.Marshal(data)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))

	// Retrieve the ID token from the response
	var bodyData map[string]interface{}
	json.NewDecoder(res.Body).Decode(&bodyData)
	token := bodyData["idToken"].(string)

	// Create a test Request with the ID token and a ResponseWriter
	r, err := http.NewRequest("", "", nil)
	addTokenToRequest(r, testToken(token))
	w := httptest.NewRecorder()

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// Call AuthenticateRequest with the test data
	handler.ServeHTTP(w, r)

	var correctStatus int
	if shouldSucceed {
		correctStatus = http.StatusOK
	} else {
		correctStatus = http.StatusUnauthorized
	}

	if w.Code != correctStatus {
		t.Fatalf("AuthorizationMiddleware test failed: response code was %v, should have been %v", w.Code, correctStatus)
	}
}

// TestAuthenticationMiddlewareValidToken tests the AuthenticateMiddleware with a valid ID Token and a dummy terminal handler
func TestAuthenticationMiddlewareValidToken(t *testing.T) {
	testMiddlewareHelper(t, true, func(validToken string) string {
		return validToken
	})
}

func TestAuthenticationMiddleWareNoToken(t *testing.T) {
	testMiddlewareHelper(t, false, func(validToken string) string {
		return ""
	})
}

func TestAuthenticationMiddleWareInvalidToken(t *testing.T) {
	testMiddlewareHelper(t, false, func(validToken string) string {
		return "INVALID_TOKEN"
	})
}
