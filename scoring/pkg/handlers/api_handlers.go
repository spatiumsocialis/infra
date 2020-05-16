// Package handlers contains the HTTP handlers for the circle service
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/models"
)

// GetCircleScoreForPeriod handles requests to get a circle score for a period
func GetCircleScoreForPeriod(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Get request vars
		vars := mux.Vars(r)

		// Get period
		period := vars[config.PeriodParameterString]
		var circleScore models.CircleScore

		if period == "" || period == "2week" {
			circleScore = models.GetRollingWindowCircleScoreForDate(user, time.Now(), s.DB)
		} else if period == "day" {
			circleScore = models.GetDayCircleScoreForDate(user, time.Now(), s.DB)
		} else {
			err := fmt.Errorf("Error: '%s' is not a valid period", period)
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Marshal payload
		payload, err := json.Marshal(circleScore)
		if err != nil {
			log.Fatalf("Error marshalling circleScore: %v", err.Error())
		}

		// Write to the response
		w.Write(payload)
	})
}

// GetEventScoresForPeriod returns event scores in a given period
func GetEventScoresForPeriod(s *common.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Get the current user
		user, err := auth.GetUser(r, s.DB)
		if err != nil {
			http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Get request vars
		vars := mux.Vars(r)

		// Get period
		period := vars[config.PeriodParameterString]
		var eventScores []models.EventScore

		switch period {
		case "all":
			if err := s.DB.Find(&eventScores).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "2week":
			eventScores, err = models.GetEventsInRollingWindow(s.DB, user, time.Now())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "day":
			eventScores, err = models.GetEventsOnDay(s.DB, user, time.Now())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			err := fmt.Errorf("Error: '%s' is not a valid period", period)
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Marshal payload
		payload, err := json.Marshal(eventScores)
		if err != nil {
			log.Fatalf("Error marshalling eventScores: %v", err.Error())
		}

		// Write to the response
		w.Write(payload)
	})
}
