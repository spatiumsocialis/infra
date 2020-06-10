package routes

import (
	"fmt"
	"strings"

	"github.com/spatiumsocialis/infra/configs/services/scoring/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/services/scoring/handlers"
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
