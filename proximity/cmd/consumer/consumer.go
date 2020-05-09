package main

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
)

func handleInteractionMesssage(m *sarama.ConsumerMessage) error {
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
	common.NewConsumer(handleInteractionMesssage)
}
