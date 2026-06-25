package declarations

import (
	"github.com/smtdfc/nagare/core/tool"
)

func InitTools() {
	tool.GlobalToolRegistry.Register(geolocation)
	tool.GlobalToolRegistry.Register(get_weather)
	tool.GlobalToolRegistry.Register(get_local_time)
	tool.GlobalToolRegistry.Register(search_wikipedia)
	tool.GlobalToolRegistry.Register(get_wikipedia_page)
	tool.GlobalToolRegistry.Register(search_github)
	tool.GlobalToolRegistry.Register(open_web)
	tool.GlobalToolRegistry.Register(get_system_info)
	tool.GlobalToolRegistry.Register(power_control)
}
