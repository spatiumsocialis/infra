package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Interaction{},
	&auth.User{},
}

// Interaction represents a proximity interaction event between a user and another person
type Interaction struct {
	gorm.Model
	// Firebase UID of the user
	UID string `json:"uid"`
	// Average distance between the user and the other party in centimeters
	Distance float32 `json:"distance"`
	// Duration of the interaction in nanoseconds
	Duration time.Duration `json:"duration"`
	// Timestamp of the beginning of the interaction
	Timestamp time.Time `json:"timestamp"`
	// Score for the interaction
	Score float32 `json:"score"`
}
