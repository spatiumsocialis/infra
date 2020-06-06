package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/common/auth"
)

// SaveUpdatedUserMessageHandler decodes the user from the message and saves it
func SaveUpdatedUserMessageHandler(s *common.Service, m *sarama.ConsumerMessage) error {
	var ole ObjectLogEntry

	if err := json.Unmarshal(m.Value, &ole); err != nil {
		return fmt.Errorf("error unmarshalling user message: %v", err)
	}

	var user auth.User
	if err := json.Unmarshal(ole.Object, &user); err != nil {
		return fmt.Errorf("error decoding user message: %v", err)
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}

	return nil

}
