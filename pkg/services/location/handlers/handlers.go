package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/services/location/models"
)

// AddPing handles requests for adding new pings
func AddPing(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Decode the interaction from the request body
		var ping models.Ping
		if err := json.NewDecoder(r.Body).Decode(&ping); err != nil {
			http.Error(w, "Error decoding ping from request: "+err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("request ping: %+v\n", ping)
		s.DB.Create(&ping)
		log.Printf("created ping: %+v\n", ping)
		json.NewEncoder(w).Encode(&ping)

		// // Log a new interaction (send msg to kafka)
		// common.LogObject(s.Producer, string(ping.ID), ping, s.ProductionTopic)

	})
}

// GetPings handles requests to get all the pings
func GetPings(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var pings []models.Ping
		if err := s.DB.Find(&pings).Error; err != nil {
			common.ThrowError(w, fmt.Errorf("error finding pings: %v", err), http.StatusInternalServerError)

		}
		if err := json.NewEncoder(w).Encode(&pings); err != nil {
			common.ThrowError(w, fmt.Errorf("error encoding pings: %v", err), http.StatusInternalServerError)

		}
	})
}
