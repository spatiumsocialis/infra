package common

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

const (
	ProductionTopic = "interaction_added"
	TestVal         = 5
)

var (
	s *Service
	p sarama.AsyncProducer
)

type testObject struct {
	val int
}

func TestMain(m *testing.M) {
	flag.Parse()

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if brokerList == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokers := strings.Split(brokerList, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokers, ", "))

	p = NewObjectLogProducer()

	os.Exit(m.Run())

}

// TestInteractionHandler tests the interaction handler by sending a POST request to interactionHandler
// to create a new Interaction, followed by a GET request to retrieve it, and ensuring  the two results are the same.
func TestLogObject(t *testing.T) {
	// Assign mock producer to service

	go func() {
		for success := range p.Successes() {
			t.Logf("producer success: %+v\n", success.Value)
		}
	}()

	// Create a test interaction and a test token
	o := &testObject{val: TestVal}

	LogObject(p, string(o.val), o, ProductionTopic)

	time.Sleep(5 * time.Second)
}
