package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spatiumsocialis/infra/configs/services/location/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/auth"
	"github.com/spatiumsocialis/infra/pkg/services/location/models"
)

var s *common.Service

var saramaConfig *sarama.Config

func TestMain(m *testing.M) {
	if err := common.LoadEnv(); err != nil {
		log.Fatalln(err)
	}
	os.Setenv("DB_PROVIDER", "sqlite3")
	os.Setenv("DB_CONNECTION_STRING", ":memory:")

	saramaConfig = sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	s = common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, nil, config.ProductionTopic)
	os.Exit(m.Run())

}

// TestAddPingHandler tests the AddPing handler
func TestAddPingHandler(t *testing.T) {
	// Assign mock producer to service
	mp := mocks.NewAsyncProducer(t, saramaConfig)

	go func() {
		for success := range mp.Successes() {
			t.Logf("producer success: %+v\n", success.Value)
		}
	}()

	s.Producer = mp
	// Create a test interaction and a test token
	testPing := models.Ping{Lat: 45.081145, Lon: -72.287415, Timestamp: time.Now()}
	testUID := "TEST_UID"
	testToken := &auth.Token{UID: testUID}
	// Marshal the text interaction to JSON, as it would be received in a POST request
	payload, err := json.Marshal(testPing)
	if err != nil {
		t.Fatalf("Error marshaling test ping: %v", err.Error())
	}

	// Create a test request and add the test token to its context
	r := httptest.NewRequest("POST", "/location", bytes.NewBuffer(payload))
	ctx := auth.AddTokenTo(context.Background(), testToken)
	w := httptest.NewRecorder()
	// Call the interaction handler with the response recorder and test request
	mp.ExpectInputAndSucceed()
	AddPing(s).ServeHTTP(w, r.WithContext(ctx))

	if w.Code != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Result().Body)
		t.Fatalf("Error adding ping: %v", string(body))
	}
	//  Read the body of the response recorder
	resBuffer := bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var createPing models.Ping
	err = json.Unmarshal(resBuffer.Bytes(), &createPing)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Ping")
	}

	t.Logf("POST response ping: %+v", createPing)

	// Make a GET request to retrieve the interaction
	r = httptest.NewRequest("GET", "/location", nil)
	w = httptest.NewRecorder()
	// Call the interaction handler with the response recorder and test request
	mp.ExpectInputAndSucceed()
	GetPings(s).ServeHTTP(w, r.WithContext(ctx))

	if w.Code != http.StatusOK {
		body, _ := ioutil.ReadAll(w.Result().Body)
		t.Fatalf("Error adding ping: %v", string(body))
	}
	//  Read the body of the response recorder
	resBuffer = bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("error reading from response buffer: %v", err)
	}

	// Unmarshal the returned interaction
	var getPings []models.Ping
	err = json.Unmarshal(resBuffer.Bytes(), &getPings)
	if err != nil {
		t.Fatalf("error unmarshalling response body into Ping: %v\n", err)
	}
	if len(getPings) != 1 {
		t.Fatalf("Fail: expecting 1 ping, got %v", len(getPings))
	}

	// Retrieve the sole ping
	getPing := getPings[0]

	t.Logf("GET response ping: %+v", getPing)

	// Check that the two pings are equal
	if createPing != getPing {
		t.Fatal("Fail: ping returned by POST and by GET are not equal")
	}
}
