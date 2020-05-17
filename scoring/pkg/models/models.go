package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/scoring/config"
)

// Schema holds the models to be included in the db schema
var Schema = common.Schema{
	&EventScore{},
	&auth.User{},
}

// EventType is an enum type for representing the different types of scored events, such as proximity interactions and daily rewards
type EventType uint

const (
	// ProximityInteraction is the type for proximity interaction events
	ProximityInteraction EventType = iota
	// DailyAllowance is the points users get once a day
	DailyAllowance
)

var eventTypeToString = map[EventType]string{
	ProximityInteraction: config.ProximityInteractionEventTypeString,
	DailyAllowance:       config.DailyAllowanceEventTypeString,
}

func (e EventType) String() string {
	return eventTypeToString[e]
}

// EventScore represents the scoring of an event
type EventScore struct {
	gorm.Model `json:"-"`
	UID        string    `json:"uid"`
	EventID    uint      `json:"eventId"`
	EventType  EventType `json:"eventType"`
	Timestamp  time.Time `json:"timestamp"`
	Score      int       `json:"score"`
}

type eventScoreResponse struct {
	UID       string
	EventID   uint
	EventType string
	Timestamp time.Time
	Score     int
}

// MarshalJSON returns a marshalled EventScore
func (e EventScore) MarshalJSON() ([]byte, error) {
	r := eventScoreResponse{
		UID:       e.UID,
		EventID:   e.EventID,
		EventType: e.EventType.String(),
		Timestamp: e.Timestamp,
		Score:     e.Score,
	}
	return json.Marshal(r)
}

// CreateEventScore creates a new EventScore object and writes it to the database before returning it
func CreateEventScore(db *gorm.DB, uid string, eventID uint, eventType EventType, timestamp time.Time, score int) (*EventScore, error) {
	es := EventScore{
		UID:       uid,
		EventID:   eventID,
		EventType: eventType,
		Timestamp: timestamp,
		Score:     score,
	}
	if err := db.Create(&es).Error; err != nil {
		return &es, err
	}
	log.Printf("event score created: %+v\n", es)
	return &es, nil
}

// UserScore is a mapping between a user and their score
type UserScore struct {
	UID   string
	Score int
}

// CircleScore represents the total and individual scores for a circle
type CircleScore struct {
	CircleID   string
	Score      int
	UserScores []UserScore
}

func calculateCircleScore(user auth.User, scores []EventScore) CircleScore {
	userScores := make([]UserScore, len(scores))
	total := 0
	for i, s := range scores {
		userScores[i] = UserScore{UID: s.UID, Score: s.Score}
		total += s.Score
	}
	return CircleScore{CircleID: user.CircleID, Score: total, UserScores: userScores}
}

// GetEventsInRange returns all the user's events in the given date range
func getEventsInRange(db *gorm.DB, user auth.User, start time.Time, end time.Time) ([]EventScore, error) {
	var eventScores []EventScore
	if err := db.Where("uid = ? OR uid = ? AND timestamp BETWEEN ? AND ?", user.ID, config.AllUserID, start.Format(time.RFC3339), end.Format(time.RFC3339)).
		Find(&eventScores).
		Error; err != nil {
		return eventScores, err
	}
	return eventScores, nil
}

// GetEventsInRollingWindow returns the events that occured in the rolling window preceding date
func GetEventsInRollingWindow(db *gorm.DB, user auth.User, date time.Time) ([]EventScore, error) {
	end := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	start := end.AddDate(0, 0, -config.RollingWindowDays)
	return getEventsInRange(db, user, start, end)
}

// GetEventsOnDay returns the events that occured on the given date
func GetEventsOnDay(db *gorm.DB, user auth.User, date time.Time) ([]EventScore, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	return getEventsInRange(db, user, start, end)
}

func getCircleScoreForDates(user auth.User, startDate time.Time, endDate time.Time, db *gorm.DB) CircleScore {
	scores := make([]EventScore, 0)
	if user.CircleID != "" {
		db.Table(
			"scores",
		).Select(
			"scores.uid",
			"users.circle_id",
			"scores.value",
		).Joins(
			"left join users on users.id = scores.uid",
		).Where(
			"users.circle_id = ? AND timestamp BETWEEN ? AND ?",
			user.CircleID,
			startDate.Format(time.RFC3339),
			endDate.Format(time.RFC3339),
		).Find(&scores)
	} else {
		db.Where("timestamp BETWEEN ? AND ? AND uid = ?", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), user.ID).Find(&scores)
	}
	return calculateCircleScore(user, scores)
}

// GetDayCircleScoreForDate calculates a user's circle's scores for a certain date
func GetDayCircleScoreForDate(user auth.User, date time.Time, db *gorm.DB) CircleScore {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	return getCircleScoreForDates(user, start, end, db)
}

// GetRollingWindowCircleScoreForDate calculates a user's circle's scores for a two week rolling window ending on date
func GetRollingWindowCircleScoreForDate(user auth.User, date time.Time, db *gorm.DB) CircleScore {
	end := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	start := end.AddDate(0, 0, -config.RollingWindowDays)
	return getCircleScoreForDates(user, start, end, db)
}
