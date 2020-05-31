package config

import (
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/pkg/auth"
)

// ServiceName is the name of the service
const ServiceName = "Proximity service"

// ServicePathPrefix is the path prefix which this service's endpoints will have
const ServicePathPrefix = "/proximity"

// Middleware holds the list of middlewares to be employed by this service
var Middleware = []mux.MiddlewareFunc{
	auth.Middleware,
}

// ProductionTopic is the Kafka topic this service produces to
const ProductionTopic = "interaction_added"

const debouncingPeriodEnvVariableName = "INTERACTION_DEBOUNCING_PERIOD_SECONDS"

// InteractionDebouncingPeriod is the number of seconds after receiving an interaction between two users for which subsequent interactions should be ignored
func InteractionDebouncingPeriod() int {
	period, err := strconv.Atoi(os.Getenv(debouncingPeriodEnvVariableName))
	if err != nil {
		return defaultInteractionDebouncingPeriod
	}
	return period
}

const defaultInteractionDebouncingPeriod = 60
