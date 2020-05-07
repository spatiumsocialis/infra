package circle

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/safe-distance/circle/pkg/handlers"
	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
)

// ServicePrefix  is the URL path prefix for the service
const ServicePrefix = "/api/v1/circle"

// Route represents one route of the service
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter returns a router for this service
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = auth.Middleware(handler)
		handler = common.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(ServicePrefix + route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = []Route{
	{
		"AddToCircle",
		strings.ToUpper("Patch"),
		"/add",
		handlers.AddToCircle,
	},

	{
		"GetCircle",
		strings.ToUpper("Get"),
		"/",
		handlers.GetCircle,
	},
	{
		"RemoveFromCircle",
		strings.ToUpper("Patch"),
		"/remove",
		handlers.RemoveFromCircle,
	},
}
