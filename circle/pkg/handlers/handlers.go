// Package handlers contains the HTTP handlers for the circle service
package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/circle/config"
	"github.com/safe-distance/socium-infra/circle/pkg/models"
	"github.com/safe-distance/socium-infra/common"
)

// AddToCircle returns a handler which adds the current user to the specified circle
func AddToCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body :"+err.Error(), http.StatusBadRequest)
			return
		}
		// Unmarshal the circle
		var circle models.Circle
		err = json.Unmarshal(body, &circle)
		if err != nil {
			http.Error(w, "Error unmarshalling circle :"+err.Error(), http.StatusBadRequest)
			return
		}

		// Start association mode
		association := s.DB.Model(&circle).Association("Users")
		if association.Error != nil {
			log.Fatalf("Error entering association mode: %v", association.Error.Error())
		}

		// Add the user to the group
		association.Append(&user)

		// Retrieve the updated circle
		s.DB.Preload("Users").Find(&circle, models.Circle{ID: user.CircleID})

		payload, err := json.Marshal(circle)
		if err != nil {
			log.Fatalf("Error marshalling circle: %v", err.Error())
		}

		// Log updated user
		common.LogObject(s.Producer, user.ID, user, config.ProductionTopic)

		// Write the user back to the response
		w.Write(payload)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}

// RemoveFromCircle removes the current user from their circle
func RemoveFromCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Check if user is in a circle
		if user.CircleID == "" {
			http.Error(w, "Error: current user isn't in a circle", http.StatusBadRequest)
			return
		}
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Unmarshal the JSON
		var bodyJSON map[string]*json.RawMessage
		json.Unmarshal(body, &bodyJSON)
		var removeUserID string
		json.Unmarshal(*(bodyJSON["uid"]), &removeUserID)

		// Get user to remove
		var userToRemove auth.User
		if result := s.DB.Find(&userToRemove, auth.User{ID: removeUserID}); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusBadRequest)
			return
		}

		// Check that userToRemove is in user's circle
		if user.CircleID != userToRemove.CircleID {
			http.Error(w, "Error: User "+removeUserID+" is not in your circle", http.StatusBadRequest)
			return
		}

		// Get current user's circle
		var circle models.Circle
		s.DB.Preload("Users").First(&circle, models.Circle{ID: user.CircleID})

		// Start association mode
		association := s.DB.Model(&circle).Association("Users")
		if association.Error != nil {
			log.Fatalf("Error entering association mode: %v", association.Error.Error())
		}

		// Remove the user from the group
		association.Delete(&userToRemove)

		// Retrieve the updated circle
		s.DB.Preload("Users").Find(&circle, models.Circle{ID: user.CircleID})

		payload, err := json.Marshal(circle)
		if err != nil {
			log.Fatalf("Error marshalling circle: %v", err.Error())
		}
		// Write the circle back to the response
		w.Write(payload)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}

// GetCircle writes the current user's circle to the response
func GetCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Fetch the current user's circle
		var circle models.Circle
		s.DB.Preload("Users").First(&circle, models.Circle{ID: user.CircleID})
		// Marshal the circle to JSON
		payload, err := json.Marshal(circle)
		if err != nil {
			log.Fatalf("Error marshalling circle: %v", err.Error())
		}
		// Write the JSON payload to the response
		w.Write(payload)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}
