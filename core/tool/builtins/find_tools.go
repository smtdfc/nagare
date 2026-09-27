package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type FindToolsInput struct {
	Keywords []string `json:"keywords" jsonschema_description:"A list of keywords to search for relevant tools, resources, and plugins. The search will return tools that match these keywords, helping to discover available capabilities."`
}

type FindToolsOutput struct {
	Tools []tool.ToolMetadata `json:"tools,omitempty"`
}

var FindTools = tool.DefineTool(
	"find_tools",
	"Searches for tools based on provided keywords. This tool helps discover available capabilities that match the specified search criteria.",
	func(ctx *context.ExecuteContext, args *FindToolsInput, bindings tool.Bindings) (*FindToolsOutput, error) {
		foundTools, err := bindings.FindToolsByKeywords(ctx, args.Keywords)
		if err != nil {
			return nil, err
		}

		fmt.Println("Searching for tools with keywords:", args.Keywords)
		fmt.Println("Found tools:", foundTools)

		return &FindToolsOutput{
			Tools: foundTools,
		}, nil
	},
)
