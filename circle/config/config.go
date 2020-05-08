package config

import (
	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
)

// ServiceName is the name of the service
const ServiceName = "Circle service"

// ServicePathPrefix is the path prefix which this service's endpoints will have
const ServicePathPrefix = "/circle"

// Middleware holds the list of middlewares to be employed by this service
var Middleware = []mux.MiddlewareFunc{
	auth.Middleware,
}
