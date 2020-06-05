package models

import (
	"time"

	"github.com/safe-distance/socium-infra/pkg/auth"
	"github.com/safe-distance/socium-infra/pkg/common"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Interaction{},
	&auth.User{},
}

// Interaction represents a proximity interaction event between a user and another person
type Interaction struct {
	ID uint `gorm:"primary_key" json:"id"`
	// Firebase UID of the user
	UID string `json:"uid"`
	// Firebase UID of the other user
	OtherUID string `json:"-"`
	// Average distance between the user and the other party in centimeters
	Distance float32 `json:"distance"`
	// Duration of the interaction in nanoseconds
	Duration time.Duration `json:"duration"`
	// Timestamp of the beginning of the interaction
	Timestamp time.Time `json:"timestamp"`
	// Score for the interaction
	Score     float32    `json:"score"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}
