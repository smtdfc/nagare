package tool

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool/built_in"
)

var ToolRegistry = map[string]domains.Tool{}

func RegisterTool(tool domains.Tool) {
	name := tool.GetName()
	ToolRegistry[name] = tool
}

func init() {
	RegisterTool(built_in.WeatherTool)
}
