package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
)

// AddInteraction handles requests for adding new interactions
func AddInteraction(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		// Add the user's UID from the auth token to the interaction
		interaction.UID = user.ID
		s.DB.Create(&interaction)
		json.NewEncoder(w).Encode(&interaction)

		// Log a new interaction (send msg to kafka)
		common.LogObject(s.Producer, string(interaction.ID), interaction, s.ProductionTopic)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}

// GetInteractions handles requests to get the current user's interaction
func GetInteractions(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Find the user's interactions and write them to the response
		interactions := make([]models.Interaction, 0)
		s.DB.Find(&interactions, models.Interaction{UID: user.ID})
		json.NewEncoder(w).Encode(interactions)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}
