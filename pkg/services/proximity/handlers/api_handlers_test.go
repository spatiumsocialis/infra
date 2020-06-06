package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/safe-distance/socium-infra/configs/services/proximity/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/common/auth"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/models"
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

// TestInteractionHandler tests the interaction handler by sending a POST request to interactionHandler
// to create a new Interaction, followed by a GET request to retrieve it, and ensuring  the two results are the same.
func TestInteractionHandler(t *testing.T) {
	// Assign mock producer to service
	mp := mocks.NewAsyncProducer(t, saramaConfig)

	go func() {
		for success := range mp.Successes() {
			t.Logf("producer success: %+v\n", success.Value)
		}
	}()

	s.Producer = mp
	// Create a test interaction and a test token
	testInteraction := models.Interaction{Distance: 51, Duration: time.Duration(60e9), Timestamp: time.Now()}
	testUID := "TEST_UID"
	testToken := &auth.Token{UID: testUID}
	// Marshal the text interaction to JSON, as it would be received in a POST request
	payload, err := json.Marshal(testInteraction)
	if err != nil {
		t.Fatalf("Error marshaling test interaction: %v", err.Error())
	}

	// Create a test request and add the test token to its context
	r := httptest.NewRequest("POST", "/interactions", bytes.NewBuffer(payload))
	ctx := auth.AddTokenTo(context.Background(), testToken)
	w := httptest.NewRecorder()
	// Call the interaction handler with the response recorder and test request
	mp.ExpectInputAndSucceed()
	AddInteraction(s).ServeHTTP(w, r.WithContext(ctx))

	//  Read the body of the response recorder
	resBuffer := bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var createInteraction models.Interaction
	err = json.Unmarshal(resBuffer.Bytes(), &createInteraction)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Interaction")
	}

	t.Logf("POST response interaction: %+v", createInteraction)

	// Make a GET request to retrieve the interaction
	r = httptest.NewRequest("GET", "/interactions", nil)
	w = httptest.NewRecorder()
	// Call the interaction handler with the response recorder and test request
	mp.ExpectInputAndSucceed()
	GetInteractions(s).ServeHTTP(w, r.WithContext(ctx))

	//  Read the body of the response recorder
	resBuffer = bytes.NewBuffer([]byte{})
	_, err = resBuffer.ReadFrom(w.Result().Body)
	if err != nil {
		t.Fatalf("Error reading from response buffer: %v", err.Error())
	}

	// Unmarshal the returned interaction
	var getInteractions []models.Interaction
	err = json.Unmarshal(resBuffer.Bytes(), &getInteractions)
	if err != nil {
		t.Fatalf("Error unmarshalling response body into Interaction")
	}
	if len(getInteractions) != 1 {
		t.Fatalf("Fail: expecting 1 interaction, got %v", len(getInteractions))
	}

	// Retrieve the sole interaction
	getInteraction := getInteractions[0]

	t.Logf("GET response interaction: %+v", getInteraction)

	// Check that the two interactions are equal
	if createInteraction != getInteractions[0] {
		t.Fatal("Fail: interaction returned by POST and by GET are not equal")
	}
}
