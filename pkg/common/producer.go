package common

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"
)

// verifyKafkaProducerFlags does something
func verifyKafkaProducerFlags() {
	return
}

func createTLSConfiguration() (t *tls.Config) {
	if certFile != "" && keyFile != "" && caFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}

// LogObject sends an object to the appropriate kafka topic
func LogObject(p sarama.AsyncProducer, key string, o interface{}, topic string) {
	entry := NewObjectLogEntry(o)
	p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: entry,
	}
}

// NewObjectLogProducer creates a new object log AsyncProducer
func NewObjectLogProducer() sarama.AsyncProducer {
	verifyKafkaProducerFlags()
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if brokerList == "" {
		flag.PrintDefaults()
	}

	brokers := strings.Split(brokerList, ",")
	fmt.Printf("Kafka brokers: %s\n", strings.Join(brokers, ", "))

	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	tlsConfig := createTLSConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	config.Version = sarama.V2_5_0_0
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to start Sarama producer: %v\n", err))
	}

	go func() {
		for success := range producer.Successes() {
			fmt.Println("Successfully wrote entry:", success)
		}
	}()

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			fmt.Println("Failed to write new object log entry:", err)
		}
	}()

	return producer
}

// NullAsyncProducer is a struct which implements the samara.AsyncProducer interface but does nothing
type NullAsyncProducer struct{}

// AsyncClose triggers a shutdown of the producer. The shutdown has completed
// when both the Errors and Successes channels have been closed. When calling
// AsyncClose, you *must* continue to read from those channels in order to
// drain the results of any messages in flight.
func (p NullAsyncProducer) AsyncClose() {
	return
}

// Close shuts down the producer and waits for any buffered messages to be
// flushed. You must call this function before a producer object passes out of
// scope, as it may otherwise leak memory. You must call this before calling
// Close on the underlying client.
func (p NullAsyncProducer) Close() error {
	return nil
}

// Input is the input channel for the user to write messages to that they
// wish to send.
func (p NullAsyncProducer) Input() chan<- *sarama.ProducerMessage {
	return make(chan *sarama.ProducerMessage, 10)
}

// Successes is the success output channel back to the user when Return.Successes is
// enabled. If Return.Successes is true, you MUST read from this channel or the
// Producer will deadlock. It is suggested that you send and read messages
// together in a single select statement.
func (p NullAsyncProducer) Successes() <-chan *sarama.ProducerMessage {
	return make(chan *sarama.ProducerMessage)
}

// Errors is the error output channel back to the user. You MUST read from this
// channel or the Producer will deadlock when the channel is full. Alternatively,
// you can set Producer.Return.Errors in your config to false, which prevents
// errors to be returned.
func (p NullAsyncProducer) Errors() <-chan *sarama.ProducerError {
	return make(chan *sarama.ProducerError)
}

// NewNullAsyncProducer returns a new NullAsyncProducer struct
func NewNullAsyncProducer() sarama.AsyncProducer {
	return NullAsyncProducer{}
}
