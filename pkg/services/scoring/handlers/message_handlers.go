package handlers

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/safe-distance/socium-infra/configs/services/scoring/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/common/auth"
	"github.com/safe-distance/socium-infra/pkg/common/kafka"
	"github.com/safe-distance/socium-infra/pkg/services/scoring/models"
	"github.com/safe-distance/socium-infra/pkg/services/scoring/models/messages"
)

// handleInteractionAddedMessage handles messages on the interaction_added topic (serialized proximity.Interaction objects)
func handleInteractionAddedMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole kafka.ObjectLogEntry
	json.Unmarshal(m.Value, &ole)
	var i messages.ProximityInteraction
	json.Unmarshal(ole.Object, &i)

	log.Printf("umarshalled interaction: %+v\n", i)

	// Create new event score
	models.CreateEventScore(s.DB, i.UID, i.ID, models.ProximityInteraction, i.Timestamp, config.ProximityInteractionPoints)

	// TODO: Produce message for proximity service to consume

	return nil
}

// handleDailyPointsAddedMessage handles messages on the daily_allowance_awarded topic (serialized models.EventScore objects)
func handleDailyPointsAddedMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole kafka.ObjectLogEntry

	if err := json.Unmarshal(m.Value, &ole); err != nil {
		return err
	}
	var e models.EventScore
	if err := json.Unmarshal(ole.Object, &e); err != nil {
		return err
	}

	log.Printf("unmarshalled daily points event: %+v\n", e)

	// Get all users
	var users []auth.User
	if err := s.DB.Find(&users).Error; err != nil {
		return err
	}

	if _, err := models.CreateEventScore(s.DB, e.UID, e.EventID, models.DailyAllowance, e.Timestamp, e.Score); err != nil {
		return err
	}
	return nil
}

// TopicHandlerMap maps topic names to the handlers which handle messages consumed from them
var TopicHandlerMap = map[string]kafka.MessageHandler{
	"interaction_added":        handleInteractionAddedMessage,
	config.DailyAllowanceTopic: handleDailyPointsAddedMessage,
	"user_modified":            kafka.SaveUpdatedUserMessageHandler,
}
