package tool

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool/builtins"
)

var ToolRegistry = map[string]domains.Tool{}

func RegisterTool(tool domains.Tool) {
	name := tool.GetName()
	ToolRegistry[name] = tool
}

func init() {
	RegisterTool(builtins.WeatherTool)
}
