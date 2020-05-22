package config

import (
	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
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

// InteractionDebouncingPeriod is the number of seconds after receiving an interaction between two users for which subsequent interactions should be ignored
const InteractionDebouncingPeriod = 60
