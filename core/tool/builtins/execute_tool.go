package declarations

import (
	"encoding/json"
	"errors"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type ExecuteToolInput struct {
	Name string                 `json:"name" jsonschema_description:"The exact name of the tool to be called, as discovered via 'find_tools'."`
	Args map[string]interface{} `json:"args" jsonschema_description:"A map of arguments to be passed to the specified tool. The keys should match the expected argument names for the tool, and the values should be of the appropriate types."`
}

type ExecuteToolOutput struct {
	Result *tool.Result
}

var ExecuteTool = tool.DefineTool(
	"execute_tool",
	"Call a tool dynamically by specifying the tool name and arguments.",
	func(ctx *context.ExecuteContext, args *ExecuteToolInput, bindings tool.Bindings) (*ExecuteToolOutput, error) {
		jsonStr, err := json.Marshal(args.Args)
		if err != nil {
			return nil, err
		}

		result := bindings.CallTool(ctx, args.Name, string(jsonStr))
		if !result.IsSuccess {
			return nil, errors.New(result.Result)
		}
		return &ExecuteToolOutput{
			Result: result,
		}, nil
	},
	[]string{tool.RoutingCategory},
)
