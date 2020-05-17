package models

import (
	"fmt"
	"time"

	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/circle/config"
	"github.com/safe-distance/socium-infra/common"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Circle{},
	&auth.User{},
}

// Circle represents a group of users who aren't making any attempt to social distance from each other (eg families, partners)
type Circle struct {
	ID        string      `json:"id"`
	Users     []auth.User `json:"users"`
	CreatedAt time.Time   `json:"-"`
	UpdatedAt time.Time   `json:"-"`
	DeletedAt *time.Time  `json:"-" sql:"index"`
}

// AddUserToCircle adds a user to a circle
func AddUserToCircle(s *common.Service, uid string, circleID string) (Circle, error) {
	// Find circle
	var circle Circle
	if err := s.DB.FirstOrCreate(&circle, Circle{ID: circleID}).Error; err != nil {
		return Circle{}, fmt.Errorf("Error retrieving/creating circle: %v", err.Error())
	}
	fmt.Printf("circle: %+v\n", circle)
	// Start association mode
	association := s.DB.Model(&circle).Association("Users")
	if association.Error != nil {
		return Circle{}, fmt.Errorf("Error entering association mode: %v", association.Error.Error())
	}
	user := auth.User{ID: uid}
	// Add the user to the group
	association.Append(&user)
	s.DB.Save(&user)
	// Log updated user
	common.LogObject(s.Producer, user.ID, user, config.ProductionTopic)
	return circle, nil
}
