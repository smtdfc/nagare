package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type FindToolByCategoriesInput struct {
	Categories []string `json:"categories" jsonschema:"title=Categories" jsonschema_description:"A list of categories to filter tools by. The search will return tools that belong to any of the specified categories, helping to discover relevant capabilities."`
}

type FindToolByCategoriesOutput struct {
	Tools []tool.Metadata `json:"tools"`
}

var FindToolByCategories = tool.DefineTool(
	"find_tool_by_categories",
	"Searches for tools based on provided categories. This tool helps discover available capabilities that match the specified search criteria.",
	func(ctx *context.ExecuteContext, args *FindToolByCategoriesInput, bindings tool.Bindings) (*FindToolByCategoriesOutput, error) {
		foundTools, err := bindings.FindToolsByCategories(ctx, args.Categories)
		if err != nil {
			return nil, err
		}

		fmt.Printf("found %d tools for categories %v", len(foundTools), args.Categories)

		return &FindToolByCategoriesOutput{
			Tools: foundTools,
		}, nil
	},
	[]string{tool.RoutingCategory},
)
