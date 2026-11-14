package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type FindToolByCategoriesInput struct {
	Categories []string `json:"categories" jsonschema:"title=Categories" jsonschema_description:"A list of categories to filter tools by. The search will return tools that belong to any of the specified categories, helping to discover relevant capabilities."`
}

// FoundTool is the compact view of a discovered tool. The full definition,
// including its argument schema, is promoted into the conversation context by
// the executor, so repeating it in the tool result only wastes tokens.
type FoundTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FindToolByCategoriesOutput struct {
	Tools []FoundTool `json:"tools"`
}

var FindToolByCategories = tool.DefineTool(
	"find_tool_by_categories",
	"Searches for tools based on provided categories. This tool helps discover available capabilities that match the specified search criteria.",
	func(ctx *context.ExecuteContext, args *FindToolByCategoriesInput, bindings tool.Bindings) (*FindToolByCategoriesOutput, error) {
		foundTools, err := bindings.FindToolsByCategories(ctx, args.Categories)
		if err != nil {
			return nil, err
		}

		tools := make([]FoundTool, 0, len(foundTools))
		for _, found := range foundTools {
			tools = append(tools, FoundTool{
				Name:        found.Name,
				Description: found.Description,
			})
		}

		return &FindToolByCategoriesOutput{
			Tools: tools,
		}, nil
	},
	[]string{tool.RoutingCategory},
)
