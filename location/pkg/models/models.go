package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/safe-distance/socium-infra/common"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Ping{},
}

// Ping represents a location ping
type Ping struct {
	gorm.Model `json:"-"`
	// Firebase UID of the user
	Lat float64
	// Duration of the interaction in nanoseconds
	Lon float64
	// Timestamp of the beginning of the interaction
	Timestamp time.Time
}

// PingCluster represents a cluster of location pings
type PingCluster struct {
	Count     int
	Points    []int
	Timestamp time.Time
}
