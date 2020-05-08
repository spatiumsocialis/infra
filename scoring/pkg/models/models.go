package models

import (
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
)

// EventScore represents the scoring of an event
type EventScore struct {
	gorm.Model
	UID       string
	EventType EventType
	EventID   uint
	Timestamp time.Time
	Value     int
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
	for _, score := range scores {
		userScores = append(userScores, UserScore{UID: score.UID, Score: score.Value})
		total += score.Value
	}
	return CircleScore{CircleID: user.CircleID, Score: total, UserScores: userScores}
}

func getCircleScoreForDates(user auth.User, startDate time.Time, endDate time.Time, db *gorm.DB) CircleScore {
	scores := make([]EventScore, 0)
	if user.CircleID != "" {
		db.Table(
			"scores",
		).Select(
			"scores.uid",
			"scorers.circle_id",
			"scores.value",
		).Joins(
			"left join users on users.uid = scores.uid",
		).Where(
			"users.circle_id = ? AND timestamp BETWEEN ? AND ?",
			user.CircleID,
			startDate,
			endDate,
		).Find(&scores)
	} else {
		db.Where("timestamp BETWEEN ? AND ? AND uid = ?", startDate, endDate, user.ID).Find(&scores)
	}
	return calculateCircleScore(user, scores)
}

// GetDayCircleScoreForDate calculates a user's circle's scores for a certain date
func GetDayCircleScoreForDate(user auth.User, date time.Time, db *gorm.DB) CircleScore {
	return getCircleScoreForDates(user, date, date, db)
}

// GetRollingWindowCircleScoreForDate calculates a user's circle's scores for a two week rolling window ending on date
func GetRollingWindowCircleScoreForDate(user auth.User, date time.Time, db *gorm.DB) CircleScore {
	return getCircleScoreForDates(user, date.AddDate(0, 0, -config.RollingWindowDays), date, db)
}
