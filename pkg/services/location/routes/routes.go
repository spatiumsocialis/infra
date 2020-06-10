package routes

import (
	"strings"

	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/services/location/handlers"
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
