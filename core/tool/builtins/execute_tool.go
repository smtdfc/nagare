package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteToolInput struct {
	Name string `json:"name" jsonschema_description:"The exact name of the tool to be called, as discovered via 'find_tools'."`
	Args string `json:"args" jsonschema_description:"A valid JSON string containing the arguments to be passed to the tool. You MUST provide this as a properly escaped JSON string representation (e.g., '{\"param1\": \"value1\"}'), NOT as a raw JSON object or map."`
}

type ExecuteToolOutput struct {
	Result *tool.Result
}

var ExecuteTool = tool.DefineTool(
	"execute_tool",
	"Call a tool dynamically by specifying the tool name and arguments.",
	func(ctx *context.ExecuteContext, args *ExecuteToolInput, bindings tool.Bindings) (*ExecuteToolOutput, error) {

		return &ExecuteToolOutput{
			Result: bindings.CallTool(ctx, args.Name, args.Args),
		}, nil
	},
)
