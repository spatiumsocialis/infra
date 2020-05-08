package points

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // gorm driver
	"github.com/safe-distance/auth"
	"github.com/safe-distance/initialize"
)

const (
	// PeriodParameterString is the key for the request query parameter which specifies whether score returned should be for today or overall
	PeriodParameterString = "period"
)

var db *gorm.DB

// InitializeService loads the env variables and initializes the database
func InitializeService(envFilenames ...string) {
	initialize.Env(false, envFilenames...)
	db = initialize.DB(&Score{})
}

// Handler is the HandlerFunc for the service
func Handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the auth token from the request context
	token, err := auth.GetTokenFrom(r.Context())
	if err != nil {
		log.Fatalf("Error: Request context contains no authorization token. " +
			"Did you forget to use the auth.Middleware?")
	}
	// Retrieve the Scorer based on the UID from the token
	var scorer Scorer
	db.FirstOrCreate(&scorer, Scorer{UID: token.UID})
	vars := mux.Vars(r)
	if r.Method == "" || r.Method == "GET" {
		period := vars[PeriodParameterString]
		var circleScore CircleScore
		if period == "" {
			// Get overall
			circleScore = scorer.GetOverallCircleScoresForDate(db, time.Now())
		} else if period == "today" {
			// Get day scores
			circleScore = scorer.GetDayCircleScoresForDate(db, time.Now())
		} else {
			http.Error(w, "\""+period+"\" is not a valid period", http.StatusBadRequest)
			return
		}
		payload, err := json.Marshal(circleScore)
		if err != nil {
			log.Fatalf("Error marshalling circleScore: %v", err.Error())
		}
		w.Write(payload)
	}
}
