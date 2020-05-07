package proximity

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/safe-distance/socium-infra/auth"
)

// ServicePrefix is the URL prefix for the service
const ServicePrefix = "/api/v1/proximity"

// Route represents one route of the service
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter returns a router for this service
func NewRouter() *mux.Router {
	// Create new router
	router := mux.NewRouter().StrictSlash(true)
	// Register routes
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = auth.Middleware(handler)
		handler = Logger(handler, route.Name)
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
		"Index",
		strings.ToUpper("Get"),
		"/",
		index,
	},
	{
		"AddInteraction",
		strings.ToUpper("Post"),
		"/interactions",
		AddInteraction,
	},

	{
		"GetInteractions",
		strings.ToUpper("Get"),
		"/interactions",
		GetInteractions,
	},
}
