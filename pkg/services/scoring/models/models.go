package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spatiumsocialis/infra/configs/services/scoring/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/auth"
)

type (
	// UserScore is a mapping between a user and their score
	UserScore struct {
		UID   string `json:"uid"`
		Score int    `json:"score"`
	}

	// CircleScore represents the total and individual scores for a circle
	CircleScore struct {
		CircleID   string      `json:"circleId"`
		Score      int         `json:"score"`
		UserScores []UserScore `json:"userScores"`
	}

	// EventType is an enum type for representing the different types of scored events, such as proximity interactions and daily rewards
	EventType uint

	// EventScore represents the scoring of an event
	EventScore struct {
		gorm.Model `json:"-"`
		UID        string
		EventID    uint
		EventType  EventType
		Timestamp  time.Time
		Score      int
	}

	eventScoreResponse struct {
		UID       string    `json:"uid"`
		EventID   uint      `json:"eventId"`
		EventType string    `json:"eventType"`
		Timestamp time.Time `json:"timestamp"`
		Score     int       `json:"score"`
	}

	// Period represents a period in time
	Period string
)

const (
	// ProximityInteraction is the type for proximity interaction events
	ProximityInteraction EventType = iota
	// DailyAllowance is the points users get once a day
	DailyAllowance
)

const (
	day           Period = "day"
	rollingWindow Period = "2week"
)

var (
	// Schema holds the models to be included in the db schema
	Schema = common.Schema{
		&EventScore{},
		&auth.User{},
	}

	eventTypeToString = map[EventType]string{
		ProximityInteraction: config.ProximityInteractionEventTypeString,
		DailyAllowance:       config.DailyAllowanceEventTypeString,
	}
	stringToEventType = map[string]EventType{
		config.ProximityInteractionEventTypeString: ProximityInteraction,
		config.DailyAllowanceEventTypeString:       DailyAllowance,
	}
)

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

// GetEventsInPeriod returns the events that occured on the given date
func GetEventsInPeriod(db *gorm.DB, user auth.User, date time.Time, period Period) ([]EventScore, error) {
	start, end, err := startAndEndDates(period, date)
	if err != nil {
		return []EventScore{}, err
	}
	return getEventsInRange(db, user, start, end)
}

// GetCircleScoreForDate calculates a user's circle's scores for a two week rolling window ending on date
func GetCircleScoreForDate(db *gorm.DB, user auth.User, date time.Time, period Period) (CircleScore, error) {
	start, end, err := startAndEndDates(period, date)
	if err != nil {
		return CircleScore{}, err
	}
	return getCircleScoreForDates(user, start, end, db)
}

// Valid checks whether a period is valid
func (p Period) Valid() bool {
	if p == day || p == rollingWindow {
		return true
	}
	return false
}

func (e EventType) String() string {
	return eventTypeToString[e]
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

// UnmarshalJSON unmarshals the event score object
func (e *EventScore) UnmarshalJSON(b []byte) error {
	var r eventScoreResponse
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	*e = EventScore{
		UID:       r.UID,
		EventID:   r.EventID,
		EventType: stringToEventType[r.EventType],
		Timestamp: r.Timestamp,
		Score:     r.Score,
	}
	return nil
}

func startAndEndDates(p Period, date time.Time) (start time.Time, end time.Time, err error) {
	switch p {
	case day:
		start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 0, 1)
		return
	case rollingWindow:
		end = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
		start = end.AddDate(0, 0, -config.RollingWindowDays)
		return
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("'%v' is not a valid period", p)
	}
}

// GetEventsInRange returns all the user's events in the given date range
func getEventsInRange(db *gorm.DB, user auth.User, start time.Time, end time.Time) ([]EventScore, error) {
	var eventScores []EventScore
	if err := db.Order("timestamp desc").Where("uid = ? OR uid = ? AND timestamp BETWEEN ? AND ?", user.ID, config.AllUserID, start.Format(time.RFC3339), end.Format(time.RFC3339)).
		Find(&eventScores).
		Error; err != nil {
		return eventScores, err
	}
	return eventScores, nil
}

func aggregateUserScores(user auth.User, scores []EventScore) CircleScore {
	users := make(map[string]UserScore)
	total := 0
	for _, s := range scores {
		total += s.Score
		// If event is attributed to the "All" user, skip the user aggregation bit
		if s.UID == config.AllUserID {
			continue
		}
		u := users[s.UID]
		if u.UID == "" {
			u.UID = s.UID
		}
		u.Score += s.Score
		users[s.UID] = u
	}
	userScores := make([]UserScore, 0)
	for _, user := range users {
		userScores = append(userScores, user)
	}
	return CircleScore{CircleID: user.CircleID, Score: total, UserScores: userScores}
}

func getCircleScoreForDates(user auth.User, startDate time.Time, endDate time.Time, db *gorm.DB) (CircleScore, error) {

	// Get all users in the circle
	var users []auth.User
	// If user's in a circle, get the rest of their circle
	if user.CircleID != "" {
		db.Find(&users, auth.User{CircleID: user.CircleID})
	} else {
		users = append(users, user)
	}

	uids := make([]string, len(users)+1)

	for i, user := range users {
		uids[i] = user.ID
	}

	uids[len(uids)-1] = config.AllUserID

	var eventScores []EventScore
	if err := db.Where("uid IN (?) AND timestamp BETWEEN ? AND ?", uids, startDate, endDate).Find(&eventScores).Error; err != nil {
		return CircleScore{}, err
	}

	// Aggregate circle scores
	circleScore := aggregateUserScores(user, eventScores)

	// Add EventScores with score=0 for users with no events in period
	for _, uid := range uids[:len(uids)-1] {
		found := false
		for _, userScore := range circleScore.UserScores {
			if uid == userScore.UID {
				found = true
			}
		}
		if !found {
			circleScore.UserScores = append(circleScore.UserScores, UserScore{UID: uid, Score: 0})
		}
	}

	return circleScore, nil
}
