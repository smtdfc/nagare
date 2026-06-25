package declarations

import (
	"github.com/smtdfc/nagare/core/tool"
)

var Tools = tool.ListTool{geolocation, get_weather, get_local_time, search_wikipedia, get_wikipedia_page, search_github, open_web}
