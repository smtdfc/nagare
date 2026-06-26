package declarations

import (
	"github.com/smtdfc/nagare/core/tool"
)

func InitTools(toolReg *tool.ToolRegistry) {
	toolReg.Register(get_tool_by_categories)
	toolReg.Register(geolocation)
	toolReg.Register(get_weather)
	toolReg.Register(get_local_time)
	toolReg.Register(search_wikipedia)
	toolReg.Register(get_wikipedia_page)
	toolReg.Register(search_github)
	toolReg.Register(open_web)
	toolReg.Register(get_system_info)
	toolReg.Register(power_control)
	toolReg.Register(set_volume)
	toolReg.Register(open_file_manager)
}
