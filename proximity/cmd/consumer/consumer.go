package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/config"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
	"github.com/safe-distance/socium-infra/proximity/pkg/producer"
)

var (
	addr      = flag.String("addr", ":8080", "The address to bind to")
	brokers   = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	verbose   = flag.Bool("verbose", false, "Turn on Sarama logging")
	certFile  = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile   = flag.String("key", "", "The optional key file for client authentication")
	caFile    = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl = flag.Bool("verify", false, "Optional verify ssl certificates chain")
)

func handleInteractionMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole common.ObjectLogEntry
	json.Unmarshal(m.Value, &ole)
	var interaction models.Interaction
	mapstructure.Decode(ole.Object, &interaction)

	log.Printf("umarshalled interaction: %+v\n", interaction)

	return nil
}

func main() {
	common.RegisterConsumerFlags()
	common.ParseFlags()
	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	producer := producer.NewInteractionLogProducer(brokerList)

	common.LoadEnv(false)

	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, &producer, config.ProductionTopic)
	common.NewConsumer(s, handleInteractionMessage)
}
