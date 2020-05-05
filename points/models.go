package points

import (
	"time"

	"github.com/jinzhu/gorm"
)

// eventType is an enum type for representing the different types of scored events, such as proximity interactions and daily rewards
type eventType uint

const rollingWindowDays = 14

// Score represents the scoring of an event
type Score struct {
	gorm.Model
	UID       string
	EventType eventType
	EventID   uint
	Timestamp time.Time
	Value     int
}

// Scorer represents a user for storing user-circle mappings
type Scorer struct {
	gorm.Model
	UID      string
	CircleID string
}

// CircleScore represents the total and individual scores for a circle
type CircleScore struct {
	CircleID   string
	Total      int
	userScores map[string]int
}

// Get returns the given user's score
func (c CircleScore) Get(uid string) int {
	return c.userScores[uid]
}

func calculateCircleScore(circleID string, scores []Score) CircleScore {
	userScores := make(map[string]int)
	total := 0
	for _, score := range scores {
		userScores[score.UID] += score.Value
		total += score.Value
	}
	return CircleScore{CircleID: circleID, userScores: userScores}
}

func (s Scorer) getCircleScoreForDates(db *gorm.DB, startDate time.Time, endDate time.Time) CircleScore {
	scores := make([]Score, 0)
	if s.CircleID != "" {
		db.Table(
			"scores",
		).Select(
			"scores.uid",
			"scorers.circle_id",
			"scores.value",
		).Joins(
			"left join scorers on scorers.uid = scores.uid",
		).Where(
			"scorers.circle_id = ? AND timestamp BETWEEN ? AND ?",
			s.CircleID,
			startDate,
			endDate,
		).Find(&scores)
	} else {
		db.Where("timestamp BETWEEN ? AND ? AND uid = ?", startDate, endDate, s.UID).Find(&scores)
	}
	return calculateCircleScore(s.CircleID, scores)
}

// GetDayCircleScoresForDate calculates a user's circle's scores for a certain date
func (s Scorer) GetDayCircleScoresForDate(db *gorm.DB, date time.Time) CircleScore {
	return s.getCircleScoreForDates(db, date, date)
}

// GetOverallCircleScoresForDate calculates a user's circle's scores for a two week rolling window ending on date
func (s Scorer) GetOverallCircleScoresForDate(db *gorm.DB, date time.Time) CircleScore {
	return s.getCircleScoreForDates(db, date.AddDate(0, 0, -rollingWindowDays), date)
}
