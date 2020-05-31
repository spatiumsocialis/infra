package routes

import (
	"fmt"
	"strings"

	"github.com/safe-distance/socium-infra/configs/services/scoring/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/services/scoring/handlers"
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
