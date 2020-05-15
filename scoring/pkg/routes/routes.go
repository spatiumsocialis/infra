package routes

import (
	"fmt"
	"strings"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/handlers"
)

// Routes holds the service's HTTP routes
var Routes = common.Routes{
	{
		Name:           "GetCircleScoreForPeriod",
		Method:         strings.ToUpper("Get"),
		Pattern:        fmt.Sprintf("/{%v}", config.PeriodParameterString),
		ServiceHandler: handlers.GetCircleScoreForPeriod,
	},
	{
		Name:           "GetEventScoresForPeriod",
		Method:         strings.ToUpper("Get"),
		Pattern:        fmt.Sprintf("/events/{%v}", config.PeriodParameterString),
		ServiceHandler: handlers.GetEventScoresForPeriod,
	},
}
