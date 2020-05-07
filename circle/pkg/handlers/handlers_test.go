package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/safe-distance/auth"
	"github.com/safe-distance/circle/pkg/common"
	"github.com/safe-distance/circle/pkg/models"
)

func TestMain(m *testing.M) {
	common.InitializeService(true, "../test.env")
	os.Exit(m.Run())
}

func TestUserHandler(t *testing.T) {
	// Create a test user and a test token
	testUID := "TEST_UID"
	testCircleID := "TEST_CIRCLE_ID"
	testCircle := models.Circle{ID: testCircleID}
	testToken := &auth.Token{UID: testUID}

	// Marshal the text interaction to JSON, as it would be received in a POST request
	payload, err := json.Marshal(testCircle)
	if err != nil {
		t.Fatalf("Error marshaling test circle: %v", err.Error())
	}

	// Create a test request and add the test token to its context
	r := httptest.NewRequest("PATCH", "/circles/add", bytes.NewBuffer(payload))
	ctx := auth.AddTokenTo(context.Background(), testToken)
	w := httptest.NewRecorder()
	// Call the interaction handler with the response recorder and test request
	AddToCircle(w, r.WithContext(ctx))

	if w.Code != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Result().Body)
		t.Fatalf("Error adding user to circle: %v", string(body))
	}

	//  Read the body of the response recorder
	resBuffer := bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var createdCircle models.Circle
	err = json.Unmarshal(resBuffer.Bytes(), &createdCircle)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Circle: %v", err.Error())
	}

	t.Logf("PATCH response circle: %+v", createdCircle)

	// PATCH
	newUserID := "NEW_USER_ID"
	newTestToken := &auth.Token{UID: newUserID}

	// Marshal the text interaction to JSON, as it would be received in a PATCH request
	payload, err = json.Marshal(testCircle)
	if err != nil {
		t.Fatalf("Error marshaling test circle: %v", err.Error())
	}

	// Make a PATCH request to update the user
	r = httptest.NewRequest("PATCH", "/circles", bytes.NewBuffer(payload))
	w = httptest.NewRecorder()
	ctx = auth.AddTokenTo(context.Background(), newTestToken)

	AddToCircle(w, r.WithContext(ctx))

	//  Read the body of the response recorder
	resBuffer = bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var updatedCircle models.Circle
	err = json.Unmarshal(resBuffer.Bytes(), &updatedCircle)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Circle: %v", err.Error())
	}

	t.Logf("PATCH response circle: %+v", updatedCircle)

	// Make a GET request to retrieve the interaction
	r = httptest.NewRequest("GET", "/circles", nil)
	w = httptest.NewRecorder()
	ctx = auth.AddTokenTo(r.Context(), testToken)
	// Call the interaction handler with the response recorder and test request
	GetCircle(w, r.WithContext(ctx))

	//  Read the body of the response recorder
	resBuffer = bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var retrievedCircle models.Circle
	err = json.Unmarshal(resBuffer.Bytes(), &retrievedCircle)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Circle: %v", err.Error())
	}

	t.Logf("GET response circle: %+v", retrievedCircle)

	// Check that the two interactions are equal
	if len(updatedCircle.Users) != len(retrievedCircle.Users) {
		t.Fatal("Fail: circles returned by PATCH and by GET are not equal")
	}

}
