package routes

import (
	"strings"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/pkg/handlers"
)

// Routes holds the service's routes
var Routes = common.Routes{
	{
		Name:           "AddInteraction",
		Method:         strings.ToUpper("Post"),
		Pattern:        "/interactions",
		ServiceHandler: handlers.AddInteraction,
	},

	{
		Name:           "GetInteractions",
		Method:         strings.ToUpper("Get"),
		Pattern:        "/interactions",
		ServiceHandler: handlers.GetInteractions,
	},
}
