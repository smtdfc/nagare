package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type DynamicTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        string `json:"args"`
	PluginID    string `json:"plugin_id"`
}

type SearchResourceInput struct {
	Keywords []string `json:"keywords" jsonschema_description:"A list of keywords to search for relevant tools, resources, and plugins. The search will return tools that match these keywords, helping to discover available capabilities."`
}

type SearchResourceOutput struct {
	Tools []DynamicTool `json:"tools,omitempty"`
}

var SearchResourceTool = tool.DefineTool(
	"search_resource_tool",
	"MANDATORY DISCOVERY TOOL: You MUST use this tool immediately whenever you lack the necessary tools, encounter missing capabilities, or need to explore a new domain requested by the user. "+
		"Search using precise keywords to retrieve available tools, resources, and plugins. "+
		"STRICTLY FORBIDDEN to guess, hallucinate, or misuse existing unrelated tools when a required capability is missing; you MUST invoke this discovery tool first to find the correct functions.",
	func(ctx *context.ExecuteContext, args *SearchResourceInput, bindings tool.Bindings) (*SearchResourceOutput, error) {
		dynamicTools, err := bindings.FindDynamicToolsByKeywords(ctx, args.Keywords)
		if err != nil {
			return nil, err
		}

		// Convert dynamic tools to the output format
		outputTools := make([]DynamicTool, len(dynamicTools))

		for i, dt := range dynamicTools {
			fmt.Println("found dynamic tool", dt.Name)
			outputTools[i] = DynamicTool{
				Name:        dt.Name,
				Description: dt.Description,
				Args:        dt.Args,
				PluginID:    dt.PluginID.String(),
			}
		}

		return &SearchResourceOutput{Tools: outputTools}, nil
	},
)
