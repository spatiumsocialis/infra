package routes

import (
	"strings"

	"github.com/safe-distance/socium-infra/circle/pkg/handlers"
	"github.com/safe-distance/socium-infra/common"
)

// Routes holds the service's HTTP routes
var Routes = common.Routes{
	{
		Name:           "AddToCircle",
		Method:         strings.ToUpper("Patch"),
		Pattern:        "/circles/add",
		ServiceHandler: handlers.AddToCircle,
	},

	{
		Name:           "GetCircle",
		Method:         strings.ToUpper("Get"),
		Pattern:        "/circles",
		ServiceHandler: handlers.GetCircle,
	},
	{
		Name:           "RemoveFromCircle",
		Method:         strings.ToUpper("Patch"),
		Pattern:        "/circles/remove",
		ServiceHandler: handlers.RemoveFromCircle,
	},
}
