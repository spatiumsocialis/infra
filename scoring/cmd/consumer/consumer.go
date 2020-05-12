package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/common"
	pmodels "github.com/safe-distance/socium-infra/proximity/pkg/models"
	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/models"
)

func handleInteractionMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole common.ObjectLogEntry
	json.Unmarshal(m.Value, &ole)
	var i pmodels.Interaction
	mapstructure.Decode(ole.Object, &i)

	log.Printf("umarshalled interaction: %+v\n", i)

	// Create new event score
	models.CreateEventScore(s.DB, i.UID, i.ID, models.ProximityInteraction, time.Now(), config.ProximityInteractionPoints)

	// TODO: Produce message for proximity service to consume

	return nil
}

func main() {
	common.RegisterConsumerFlags()
	common.ParseFlags()
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, , config.ProductionTopic)
	common.NewConsumer(handleInteractionMessage)
}
