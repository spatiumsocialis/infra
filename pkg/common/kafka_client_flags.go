package common

import (
	"flag"
	"os"
)

// Sarama configuration options
var (
	// Kafka bootstrap brokers to connect to, as a comma separated list
	brokerList = ""
	// Kafka cluster version
	version = "2.5.0"
	// Kafka consumer group definition
	group = ""
	// Kafka topics to be consumed, as a comma separated list
	topics = ""
	// Consumer group partition assignment strategy (range, roundrobin, sticky)
	assignor = "range"
	//  Kafka consumer consume initial offset from oldest
	oldest = true
	// Sarama logging
	verbose = false

	addr     string
	certFile string

	keyFile string

	caFile string

	verifySsl bool
)

// RegisterKafkaClientFlags registers flags needed for kafka client
func RegisterKafkaClientFlags() {
	flag.StringVar(&addr, "addr", ":8080", "The address to bind to")
	flag.StringVar(&brokerList, "brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	flag.StringVar(&certFile, "certificate", "", "The optional certificate file for client authentication")
	flag.StringVar(&keyFile, "key", "", "The optional key file for client authentication")

	flag.StringVar(&caFile, "ca", "", "The optional certificate authority file for TLS client authentication")

	flag.BoolVar(&verifySsl, "verify", false, "Optional verify ssl certificates chain")

	flag.StringVar(&group, "group", os.Getenv("KAFKA_CGROUP"), "Kafka consumer group definition")
	flag.StringVar(&version, "version", "2.5.0", "Kafka cluster version")
	flag.StringVar(&topics, "topics", "", "Kafka topics to be consumed, as a comma separated list")
	flag.StringVar(&assignor, "assignor", "range", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
	flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
	flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
}
