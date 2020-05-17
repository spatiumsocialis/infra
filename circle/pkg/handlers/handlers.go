// Package handlers contains the HTTP handlers for the circle service
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/circle/pkg/models"
	"github.com/safe-distance/socium-infra/common"
)

// AddToCircle returns a handler which adds the current user to the specified circle
func AddToCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("Error retrieving current user: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("Error reading request body: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		// Unmarshal the circle
		var circle models.Circle
		if err := json.Unmarshal(body, &circle); err != nil {
			common.ThrowError(w, fmt.Errorf("Error unmarshalling circle: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		if circle.ID == "" {
			common.ThrowError(w, fmt.Errorf("bad request: 'id' parameter missing from request body"), http.StatusBadRequest)
			return
		}

		err = models.AddUserToCircle(s, &user, &circle, true)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("error adding user to circle: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(circle)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("Error marshalling circle: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		// Write the user back to the response
		w.Write(payload)
	})
}

// RemoveFromCircle removes the current user from their circle
func RemoveFromCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		if bodyJSON["uid"] == nil {
			common.ThrowError(w, fmt.Errorf("no uid supplied"), http.StatusBadRequest)
			return
		}
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
		users := s.DB.Model(&circle).Association("Users")
		if users.Error != nil {
			log.Fatalf("Error entering association mode: %v", users.Error)
		}

		// Remove the user from the group
		users.Delete(&userToRemove)

		// Check if circle  is empty
		if users.Count() == 0 {
			if err := circle.Delete(s); err != nil {
				log.Printf("error deleting empty circle: %v", err)
			}
		}

		// Add user to a new circle
		newCircle := models.Circle{ID: uuid.New().String()}
		err = models.AddUserToCircle(s, &user, &newCircle, false)
		// save changes to user to circle
		s.DB.Find(&user)
		fmt.Printf("user: %+v circle: %+v\n", user, circle)

		payload, err := json.Marshal(circle)
		if err != nil {
			log.Fatalf("Error marshalling circle: %v", err.Error())
		}
		// Write the circle back to the response
		w.Write(payload)
	})
}

// GetCircle writes the current user's circle to the response
func GetCircle(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var circle models.Circle
		if user.CircleID == "" {
			// Generate a new circle id
			circle.ID = uuid.New().String()
			err = models.AddUserToCircle(s, &user, &circle, false)
			if err != nil {
				common.ThrowError(w, fmt.Errorf("error adding user to circle: %v", err.Error()), http.StatusInternalServerError)
				return
			}
			// Fetch the current user's circle
		} else if err := s.DB.Preload("Users").First(&circle, models.Circle{ID: user.CircleID}).Error; err != nil {
			common.ThrowError(w, fmt.Errorf("error fetching circle: %v", err), http.StatusInternalServerError)
			return
		}
		// Marshal the circle to JSON
		payload, err := json.Marshal(circle)
		if err != nil {
			log.Fatalf("Error marshalling circle: %v", err.Error())
		}
		// Write the JSON payload to the response
		w.Write(payload)
	})
}
