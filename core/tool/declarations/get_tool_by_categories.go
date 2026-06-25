package declarations

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
)

type GetToolByCategoriesArgs struct {
	Categories []string `json:"categories" jsonschema:"description=Tool groups to retrieve. You should include all potentially relevant groups (multi-group allowed) to maximize tool discovery. The backend may expand these groups to related ones automatically. Prefer broader coverage over strict minimal selection."`
}

var get_tool_by_categories = tool.DeclareTool(
	"get_tool_by_categories",
	"When no suitable tool is found in the current tool list, you MUST call this",
	func(ctx domains.AgentContext, args GetToolByCategoriesArgs) (any, error) {
		listToolAvb := tool.GlobalToolRegistry.GetToolByCategories(args.Categories)
		summarize := tool.GetToolSummary(listToolAvb)

		ctx.AfterToolCall(func(arg domains.MiddlewareArg) {
			arg.UseToolsForNextTurn(listToolAvb)
		}, true)

		return map[string]any{
			"tool_summary": summarize,
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)
