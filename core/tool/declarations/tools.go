package declarations

import (
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/utils"
	gowiki "github.com/trietmn/go-wiki"
)

func init() {
	gowiki.SetUserAgent(utils.NAGARE_USER_AGENT)
}

var appLogger = logger.GetLogger()
var Tools = tool.ListTool{geolocation, get_weather, get_local_time, search_wikipedia, get_wikipedia_page}
