package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type DynamicToolCallInput struct {
	PluginID string                 `json:"plugin_id"`
	Name     string                 `json:"name" jsonschema_description:"The name of the tool to be called. This should match the name defined in the tool's declaration."`
	Args     map[string]interface{} `json:"args" jsonschema_description:"A map of arguments to be passed to the tool. The keys should correspond to the argument names expected by the tool, and the values should be of the appropriate type."`
}

type DynamicToolCallOutput struct {
	Result any `json:"result" jsonschema_description:"The result returned by the dynamically called tool. The structure and type of this result will depend on the specific tool that was invoked."`
}

var DynamicToolCallTool = tool.DefineTool(
	"dynamic_tool_call",
	"Call a tool dynamically by specifying the plugin ID and tool name.",
	func(ctx *context.ExecuteContext, args *DynamicToolCallInput, _ tool.Bindings) (*DynamicToolCallOutput, error) {
		if args == nil || args.PluginID == "" || args.Name == "" {
			return nil, fmt.Errorf("plugin_id and name are required")
		}

		fmt.Println("called ", args.Name, args.PluginID, args.Args)
		// Call the tool dynamically using the provided plugin ID and tool name
		result := "ok"
		return &DynamicToolCallOutput{Result: result}, nil
	},
)
