// Package handlers contains the HTTP handlers for the circle service
package handlers

import (
	"encoding/json"
	"fmt"
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
			common.ThrowError(w, fmt.Errorf("error retrieving current user: %v", err), http.StatusInternalServerError)
			return
		}
		// Get request vars
		vars := mux.Vars(r)

		// Get period
		period := models.Period(vars[config.PeriodParameterString])

		if !period.Valid() {
			common.ThrowError(w, fmt.Errorf("error: '%s' is not a valid period", period), http.StatusBadRequest)
		}
		// Get circle score
		circleScore, err := models.GetCircleScoreForDate(s.DB, user, time.Now(), period)
		if err != nil {
			common.ThrowError(w, err, http.StatusInternalServerError)
		}
		// Marshal payload
		payload, err := json.Marshal(circleScore)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("error marshalling circleScore: %v", err.Error()), http.StatusInternalServerError)
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
			http.Error(w, "error retrieving current user: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Get request vars
		vars := mux.Vars(r)

		// Get period
		period := models.Period(vars[config.PeriodParameterString])

		if !period.Valid() {
			common.ThrowError(w, fmt.Errorf("error: '%s' is not a valid period", period), http.StatusBadRequest)
		}

		// Get event scores
		eventScores, err := models.GetEventsInPeriod(s.DB, user, time.Now(), period)
		if err != nil {
			common.ThrowError(w, err, http.StatusInternalServerError)
		}

		// Marshal payload
		payload, err := json.Marshal(eventScores)
		if err != nil {
			common.ThrowError(w, fmt.Errorf("error marshalling eventScores: %v", err), http.StatusInternalServerError)
		}

		// Write to the response
		w.Write(payload)
	})
}
