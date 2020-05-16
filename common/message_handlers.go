package common

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/safe-distance/socium-infra/auth"
)

// SaveUpdatedUserMessageHandler decodes the user from the message and saves it
func SaveUpdatedUserMessageHandler(s *Service, m *sarama.ConsumerMessage) error {
	var ole ObjectLogEntry

	if err := json.Unmarshal(m.Value, &ole); err != nil {
		return err
	}

	var user auth.User
	if err := mapstructure.Decode(ole.Object, &user); err != nil {
		return err
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil

}
