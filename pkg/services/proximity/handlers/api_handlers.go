package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/safe-distance/socium-infra/configs/services/proximity/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/common/auth"
	"github.com/safe-distance/socium-infra/pkg/common/kafka"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/models"
)

// AddInteraction handles requests for adding new interactions
func AddInteraction(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Decode the interaction from the request body
		var interaction models.Interaction
		if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
			http.Error(w, "Error decoding interaction from request: "+err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("request interaction: %+v\n", interaction)
		// Check whether other user is in current user's circle
		otherUserUID := interaction.UID
		var otherUser auth.User
		s.DB.First(&otherUser, auth.User{ID: otherUserUID})
		if otherUser.CircleID != "" && otherUser.CircleID == user.CircleID {
			response := make(map[string]string)
			response["error"] = "users in same circle"
			json.NewEncoder(w).Encode(response)
			return
		}
		// Check if we've already registered an interaction for these users in the debouncing window
		startOfDebouncingPeriod := interaction.Timestamp.Add(-1 * time.Duration(config.InteractionDebouncingPeriod()) * time.Second)
		fmt.Printf("start of debouncing period: %v\n", startOfDebouncingPeriod)
		var mostRecentInteraction models.Interaction
		query := s.DB.Where(models.Interaction{UID: user.ID, OtherUID: otherUserUID}).Or(models.Interaction{UID: otherUserUID, OtherUID: user.ID}).Order("timestamp desc")
		query = query.Attrs(models.Interaction{UID: "not_found"}).FirstOrInit(&mostRecentInteraction)
		if query.Error != nil {
			common.ThrowError(w, fmt.Errorf("error retrieving most recent interaction between %v and %v: %v", user.ID, otherUserUID, err), http.StatusAlreadyReported)
			return
		}
		if mostRecentInteraction.UID != "not_found" && mostRecentInteraction.Timestamp.After(startOfDebouncingPeriod) {
			msg := fmt.Sprintf("error: interaction between these two users recorded at %v", mostRecentInteraction.Timestamp)
			common.ThrowError(w, errors.New(msg), http.StatusAlreadyReported)
			return
		}

		// Add the user's UID from the auth token to the interaction
		interaction.UID = user.ID
		interaction.OtherUID = otherUserUID
		s.DB.Create(&interaction)
		json.NewEncoder(w).Encode(&interaction)
		fmt.Printf("created interaction: %+v\n", interaction)

		// Publish a new interaction (send msg to kafka)
		kafka.LogObject(s.Producer, string(interaction.ID), interaction, s.ProductionTopic)

	})
}

// GetInteractions handles requests to get the current user's interaction
func GetInteractions(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Find the user's interactions and write them to the response
		interactions := make([]models.Interaction, 0)
		s.DB.Where(models.Interaction{UID: user.ID}).Or(models.Interaction{OtherUID: user.ID}).Find(&interactions)
		// Set all interactions UID field to the current uid for the response so no one else's UID gets leaked
		for i := 0; i < len(interactions); i++ {
			interactions[i].UID = user.ID
		}
		json.NewEncoder(w).Encode(interactions)

	})
}
