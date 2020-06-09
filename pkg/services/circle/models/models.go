package models

import (
	"fmt"
	"time"

	"github.com/spatiumsocialis/infra/configs/services/circle/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/auth"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Circle{},
	&auth.User{},
}

type (
	// Circle represents a group of users who aren't making any attempt to social distance from each other (eg families, partners)
	Circle struct {
		ID        string      `json:"id"`
		Users     []auth.User `json:"users"`
		CreatedAt time.Time   `json:"-"`
		UpdatedAt time.Time   `json:"-"`
		DeletedAt *time.Time  `json:"-" sql:"index"`
	}
	// CircleResponse represents a circle as it's returned to the client, with user profiles
	CircleResponse struct {
		ID    string         `json:"id"`
		Users []auth.Profile `json:"users"`
	}
)

// NewCircleResponse returns a new response circle from the given circle
func NewCircleResponse(c Circle) (*CircleResponse, error) {
	profiles, err := auth.GetUserProfiles(c.Users...)
	if err != nil {
		return nil, err
	}
	return &CircleResponse{ID: c.ID, Users: profiles}, nil
}

// AddUserToCircle adds a user to a circle
func AddUserToCircle(s *common.Service, user *auth.User, circle *Circle, mergeCircles bool) error {
	// Find circle
	if err := s.DB.Preload("Users").FirstOrCreate(circle, Circle{ID: circle.ID}).Error; err != nil {
		return fmt.Errorf("error retrieving/creating circle: %v", err.Error())
	}
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
	if err := association.Append(user).Error; err != nil {
		fmt.Printf("error adding user to circle: %v\n", err)
		return err
	}
	if err := s.DB.Save(user).Error; err != nil {
		fmt.Printf("error saving user: %v\n", err)
		return err
	}
	if err := s.DB.Save(circle).Error; err != nil {
		fmt.Printf("error saving circle: %v\n", err)
		return err
	}
	if mergeCircles && oldCircleID != "" {
		var users []auth.User
		if err := s.DB.Find(&users, auth.User{CircleID: oldCircleID}).Error; err != nil {
			fmt.Printf("error finding other circle users: %v\n", err)
			return err
		}
		if err := association.Append(&users).Error; err != nil {
			fmt.Printf("error merging circles: %v\n", err)
			return err
		}
	}
	// Log updated user
	kafka.LogObject(s.Producer, user.ID, user, config.ProductionTopic)
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
