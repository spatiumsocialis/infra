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
	UID string
	// Average distance between the user and the other party in centimeters
	Distance float32
	// Duration of the interaction in nanoseconds
	Duration time.Duration
	// Timestamp of the beginning of the interaction
	Timestamp time.Time
	// Score for the interaction
	Score float32
}
