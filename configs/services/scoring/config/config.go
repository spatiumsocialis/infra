package config

import (
	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/pkg/common/auth"
)

// General

// ServiceName is the name of the service
const ServiceName = "Scoring service"

// ServicePathPrefix is the path prefix which this service's endpoints will have
const ServicePathPrefix = "/scores"

// Middleware holds the list of middlewares to be employed by this service
var Middleware = []mux.MiddlewareFunc{
	auth.Middleware,
}

//  Service-specific

// PeriodParameterString is the key for the 'period' path parameter
const PeriodParameterString = "period"

// RollingWindowDays is the number of days in a rolling window period
const RollingWindowDays = 14

// ProductionTopic is the Kafka topic this service produces to
const ProductionTopic = "interaction_scored"

// DailyAllowanceTopic is the Kafka topic which the daily points allowance is produced to
const DailyAllowanceTopic = "daily_allowance_awarded"

// Scores for events
const (
	ProximityInteractionPoints = -100
	DailyAllowancePoints       = 1000
)

// Event type string representations
const (
	ProximityInteractionEventTypeString = "proximity_interaction"
	DailyAllowanceEventTypeString       = "daily_allowance"
)

// AllUserID is the userID to use for event scores which apply to all users
const AllUserID = "ALL"
