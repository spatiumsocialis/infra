package common

import "github.com/Shopify/sarama"

// NullAsyncProducer is a struct which implements the samara.AsyncProducer interface but does nothing
type NullAsyncProducer struct{}

// AsyncClose triggers a shutdown of the producer. The shutdown has completed
// when both the Errors and Successes channels have been closed. When calling
// AsyncClose, you *must* continue to read from those channels in order to
// drain the results of any messages in flight.
func (p *NullAsyncProducer) AsyncClose() {
	return
}

// Close shuts down the producer and waits for any buffered messages to be
// flushed. You must call this function before a producer object passes out of
// scope, as it may otherwise leak memory. You must call this before calling
// Close on the underlying client.
func (p *NullAsyncProducer) Close() error {
	return nil
}

// Input is the input channel for the user to write messages to that they
// wish to send.
func (p *NullAsyncProducer) Input() chan<- *sarama.ProducerMessage {
	return make(chan *sarama.ProducerMessage)
}

// Successes is the success output channel back to the user when Return.Successes is
// enabled. If Return.Successes is true, you MUST read from this channel or the
// Producer will deadlock. It is suggested that you send and read messages
// together in a single select statement.
func Successes() <-chan *sarama.ProducerMessage {
	return make(chan *sarama.ProducerMessage)
}

// Errors is the error output channel back to the user. You MUST read from this
// channel or the Producer will deadlock when the channel is full. Alternatively,
// you can set Producer.Return.Errors in your config to false, which prevents
// errors to be returned.
func Errors() <-chan *sarama.ProducerError {
	return make(chan *sarama.ProducerError)
}
