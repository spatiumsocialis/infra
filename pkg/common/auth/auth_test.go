package auth

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := common.LoadEnv(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func addTokenToRequest(r *http.Request, token string) {
	r.Header.Add("Authorization", "Bearer "+token)
}

func testMiddlewareHelper(t *testing.T, shouldSucceed bool, testToken func(validToken string) string) {
	// api key retrieved from https://console.firebase.google.com/u/0/project/spatiumsocialis-e4683/settings/general
	token, err := GenerateToken(TestUID)
	assert.Nil(t, err)

	// Create a test Request with the ID token and a ResponseWriter
	r, err := http.NewRequest("", "", nil)
	assert.Nil(t, err)
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
