package routes

import (
	"strings"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/location/pkg/handlers"
)

// Routes holds the service's routes
var Routes = common.Routes{
	{
		Name:           "AddPing",
		Method:         strings.ToUpper("Post"),
		Pattern:        "/pings",
		ServiceHandler: handlers.AddPing,
	},

	{
		Name:           "GetPings",
		Method:         strings.ToUpper("Get"),
		Pattern:        "/pings",
		ServiceHandler: handlers.GetPings,
	},
}
