package common

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/auth"
)

// SaveUpdatedUserMessageHandler decodes the user from the message and saves it
func SaveUpdatedUserMessageHandler(s *Service, m *sarama.ConsumerMessage) error {
	var ole ObjectLogEntry

	if err := json.Unmarshal(m.Value, &ole); err != nil {
		return fmt.Errorf("error unmarshalling user message: %v", err)
	}

	log.Printf("ole: %+v\n", ole)

	var user auth.User
	if err := mapstructure.Decode(ole.Object, &user); err != nil {
		return fmt.Errorf("error decoding user message: %v", err)
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}

	return nil

}