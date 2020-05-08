package routes

import (
	"strings"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/scoring/pkg/handlers"
)

// Routes holds the service's HTTP routes
var Routes = common.Routes{
	{
		Name:           "GetCircleScoreForPeriod",
		Method:         strings.ToUpper("Get"),
		Pattern:        "/scores/{period}",
		ServiceHandler: handlers.GetCircleScoreForPeriod,
	},
}
