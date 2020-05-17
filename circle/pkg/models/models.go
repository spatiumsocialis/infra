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
func AddUserToCircle(s *common.Service, user *auth.User, circle *Circle, mergeCircles bool) error {
	// Find circle
	if err := s.DB.Preload("Users").FirstOrCreate(circle, Circle{ID: circle.ID}).Error; err != nil {
		return fmt.Errorf("error retrieving/creating circle: %v", err.Error())
	}
	fmt.Printf("circle: %+v\n", circle)
	oldCircleID := user.CircleID
	// Start association mode
	association := s.DB.Model(circle).Association("Users")
	if association.Error != nil {
		return fmt.Errorf("error entering association mode: %v", association.Error.Error())
	}
	for _, u := range circle.Users {
		if u.ID == user.ID {
			fmt.Println("user already in circle, no update")
			return nil
		}
	}
	// Add the user to the circle
	association.Append(user)
	s.DB.Save(user)
	s.DB.Save(circle)
	if mergeCircles && oldCircleID != "" {
		var users []auth.User
		s.DB.Find(&users, auth.User{CircleID: oldCircleID})
		association.Append(&users)
	}
	// Log updated user
	common.LogObject(s.Producer, user.ID, user, config.ProductionTopic)
	return nil
}

// Delete deletes the circle from the service's data store
func (c *Circle) Delete(s *common.Service) error {
	if s.DB.NewRecord(c) {
		return fmt.Errorf("no primary key")
	}
	if err := s.DB.Delete(c).Error; err != nil {
		return fmt.Errorf("gorm: %v", err)
	}
	return nil
}
