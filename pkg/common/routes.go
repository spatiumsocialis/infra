package common

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route represents one route of the service
type Route struct {
	Name           string
	Method         string
	Pattern        string
	ServiceHandler ServiceHandler
}

// Routes is type-alias for a Route slice
type Routes []Route

// NewRouter returns a router for this service
func NewRouter(s *Service, routes Routes, middleware ...mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.ServiceHandler(s)
		for _, mw := range middleware {
			handler = mw(handler)
		}
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(s.PathPrefix + route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
