package config

import (
	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
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
