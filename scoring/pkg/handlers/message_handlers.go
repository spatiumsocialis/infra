package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
	pmodels "github.com/safe-distance/socium-infra/proximity/pkg/models"
	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/models"
)

// handleInteractionAddedMessage handles messages on the interaction_added topic (serialized proximity.Interaction objects)
func handleInteractionAddedMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole common.ObjectLogEntry
	json.Unmarshal(m.Value, &ole)
	var i pmodels.Interaction
	json.Unmarshal(ole.Object, &i)
	// objectMap := ole.Object.(map[string]interface{})
	// timestamp, err := time.Parse(time.RFC3339, objectMap["Timestamp"].(string))
	// if err != nil {
	// 	return err
	// }
	// i := pmodels.Interaction{
	// 	UID:       objectMap["UID"].(string),
	// 	Distance:  float32(objectMap["Distance"].(float64)),
	// 	Duration:  objectMap["Duration"].(time.Duration),
	// 	Timestamp: timestamp,
	// }
	// if err := mapstructure.Decode(ole.Object, &i); err != nil {
	// 	return fmt.Errorf("error decoding interaction: %v", err)
	// }

	log.Printf("umarshalled interaction: %+v\n", i)

	// Create new event score
	models.CreateEventScore(s.DB, i.UID, i.ID, models.ProximityInteraction, time.Now(), config.ProximityInteractionPoints)

	// TODO: Produce message for proximity service to consume

	return nil
}

// handleDailyPointsAddedMessage handles messages on the daily_points_added topic (serialized models.EventScore objects)
func handleDailyPointsAddedMessage(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole common.ObjectLogEntry

	if err := json.Unmarshal(m.Value, &ole); err != nil {
		return err
	}
	var e models.EventScore
	if err := mapstructure.Decode(ole.Object, &e); err != nil {
		return err
	}

	log.Printf("unmarshalled daily points event: %+v\n", e)

	// Get all users
	var users []auth.User
	if err := s.DB.Find(&users).Error; err != nil {
		return err
	}

	if _, err := models.CreateEventScore(s.DB, config.AllUserID, e.EventID, models.DailyAllowance, time.Now(), config.DailyAllowancePoints); err != nil {
		return err
	}
	return nil
}

// TopicHandlerMap maps topic names to the handlers which handle messages consumed from them
var TopicHandlerMap = map[string]common.MessageHandler{
	"interaction_added":  handleInteractionAddedMessage,
	"daily_points_added": handleDailyPointsAddedMessage,
	"user_modified":      common.SaveUpdatedUserMessageHandler,
}
