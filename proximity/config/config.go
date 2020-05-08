package config

import (
	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
)

// ServiceName is the name of the service
const ServiceName = "Proximity service"

// ServicePathPrefix is the path prefix which this service's endpoints will have
const ServicePathPrefix = "/proximity"

// Middleware holds the list of middlewares to be employed by this service
var Middleware = []mux.MiddlewareFunc{
	auth.Middleware,
}

// Models holds the list of models that the service uses
var Models = []interface{}{
	&models.Interaction{},
	&auth.User{},
}
