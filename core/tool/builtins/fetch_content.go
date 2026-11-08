package declarations

import (
	"unicode/utf8"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/discovery"
	"github.com/smtdfc/nagare/core/tool"
)

type FetchContentInput struct {
	Provider        string `json:"provider" jsonschema_description:"The name of the information source provider to FetchContent, Available: [Wikipedia]."`
	ID              string `json:"id"`
	MaxContentBytes int    `json:"max_content_bytes"`
}

type FetchContentOutput struct {
	Result *discovery.Detail `json:"result"`
}

func truncateContent(content string, maxBytes int) string {
	if maxBytes <= 0 || len(content) <= maxBytes {
		return content
	}

	content = content[:maxBytes]

	for !utf8.ValidString(content) {
		content = content[:len(content)-1]
	}

	return content + "\n\n[Content truncated]"
}

var FetchContentTool = tool.DefineTool(
	"fetch_content_tool",
	"Fetch the full content of a previously discovered resource by its ID from a specified source provider.",
	func(
		ctx *context.ExecuteContext,
		args *FetchContentInput,
		bindings tool.Bindings,
	) (*FetchContentOutput, error) {
		result, err := bindings.GetDiscoveryMgr().Fetch(
			ctx,
			args.Provider,
			args.ID,
		)
		if err != nil {
			return nil, err
		}

		result.Content = truncateContent(result.Content, args.MaxContentBytes)
		return &FetchContentOutput{
			Result: result,
		}, nil
	},
	[]string{tool.ResearchCategory},
)
