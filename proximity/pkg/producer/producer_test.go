package producer

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/config"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
)

var s *common.Service

var p sarama.AsyncProducer

func TestMain(m *testing.M) {
	flag.Parse()

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	p = NewInteractionLogProducer(brokerList)

	os.Exit(m.Run())

}

// TestInteractionHandler tests the interaction handler by sending a POST request to interactionHandler
// to create a new Interaction, followed by a GET request to retrieve it, and ensuring  the two results are the same.
func TestLogInteraction(t *testing.T) {
	// Assign mock producer to service

	go func() {
		for success := range p.Successes() {
			t.Logf("producer success: %+v\n", success.Value)
		}
	}()

	// Create a test interaction and a test token
	testInteraction := &models.Interaction{Distance: 51, Duration: time.Duration(60e9), Timestamp: time.Now()}

	LogInteraction(p, testInteraction, config.ProductionTopic)

	time.Sleep(5 * time.Second)
}
