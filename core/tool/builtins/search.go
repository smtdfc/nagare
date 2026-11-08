package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/discovery"
	"github.com/smtdfc/nagare/core/tool"
)

type SearchInput struct {
	Provider string `json:"provider" jsonschema_description:"The name of the information source provider to search, Available: [Wikipedia]."`
	Query    string `json:"query" jsonschema_description:"The search query used to find relevant information."`
	Limit    int    `json:"limit" jsonschema_description:"The maximum number of search results to return."`
}

type SearchOutput struct {
	Results []discovery.Source `json:"results" jsonschema_description:"The information sources found by the search operation."`
}

var SearchTool = tool.DefineTool(
	"search_tool",
	"Search for information from a specified source provider and return relevant results.",
	func(
		ctx *context.ExecuteContext,
		args *SearchInput,
		bindings tool.Bindings,
	) (*SearchOutput, error) {
		sources, err := bindings.GetDiscoveryMgr().Search(
			ctx,
			args.Provider,
			args.Query,
			args.Limit,
		)
		if err != nil {
			return nil, err
		}

		return &SearchOutput{
			Results: sources,
		}, nil
	},
	[]string{tool.ResearchCategory},
)
