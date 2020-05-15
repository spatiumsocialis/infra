package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/location/pkg/models"
)

func throwError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// AddPing handles requests for adding new pings
func AddPing(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the interaction from the request body
		var ping models.Ping
		if err := json.NewDecoder(r.Body).Decode(&ping); err != nil {
			http.Error(w, "Error decoding interaction from request: "+err.Error(), http.StatusBadRequest)
			return
		}
		s.DB.Create(&ping)
		json.NewEncoder(w).Encode(&ping)

		// // Log a new interaction (send msg to kafka)
		// common.LogObject(s.Producer, string(ping.ID), ping, s.ProductionTopic)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}

// GetPings handles requests to get all the pings
func GetPings(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pings []models.Ping
		if err := s.DB.Find(&pings).Error; err != nil {
			throwError(w, fmt.Errorf("error finding pings: %v", err))

		}
		if err := json.NewEncoder(w).Encode(&pings); err != nil {
			throwError(w, fmt.Errorf("error encoding pings: %v", err))

		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}
