package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/config"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
	"github.com/safe-distance/socium-infra/proximity/pkg/producer"
	"github.com/safe-distance/socium-infra/proximity/pkg/routes"
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

func main() {
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

	producer := producer.NewInteractionLogProducer(brokerList)

	common.LoadEnv(false)
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := common.NewRouter(s, routes.Routes, config.Middleware...)
	http.Handle("/", r)
	fmt.Printf("Listening on port %v...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
