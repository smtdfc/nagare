package registry

import (
	"github.com/smtdfc/nagare/core/tool"
	declarations "github.com/smtdfc/nagare/core/tool/builtins"
)

var Registry = map[string]tool.Tool{}

func RegisterTool(tool tool.Tool) {
	Registry[tool.GetName()] = tool
}

func init() {
	RegisterTool(declarations.WeatherTool)
}
