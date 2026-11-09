package tool

import (
	"github.com/smtdfc/nagare/core/computer_use"
	context2 "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/discovery/manager"
)

type Bindings interface {
	RefreshCronJob(ctx *context2.ExecuteContext)
	CreateCronJob(ctx *context2.ExecuteContext, sessionID string, name string, prompt string, cronExpr string) (string, error)
	FindToolsByCategories(ctx *context2.ExecuteContext, categories []string) ([]Metadata, error)
	CallTool(ctx *context2.ExecuteContext, toolName string, args string) *Result
	RequestToolApproval(ctx *context2.ExecuteContext, title string) bool
	GetDiscoveryMgr() *manager.DiscoveryManager
	GetComputerUseOrchestrator() *computer_use.ComputerUseOrchestrator
}
