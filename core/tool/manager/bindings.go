package manager

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/computer_use"
	"github.com/smtdfc/nagare/core/context"
	cronjob_manager "github.com/smtdfc/nagare/core/cronjob/manager"
	"github.com/smtdfc/nagare/core/discovery/manager"
	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/hil"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/tool/registry"
)

type ToolBindings struct {
	approvalMgr             *hil.Manager
	cronJobMgr              *cronjob_manager.CronJobManager
	toolMgr                 *ToolManager
	eventBus                *event_bus.CoreEventBus
	discoveryMgr            *manager.DiscoveryManager
	computerUseOrchestrator *computer_use.ComputerUseOrchestrator
	logger                  *logger.BaseLogger
}

// GetComputerUseOrchestrator implements [tool.Bindings].
func (t *ToolBindings) GetComputerUseOrchestrator() *computer_use.ComputerUseOrchestrator {
	return t.computerUseOrchestrator
}

// GetDiscoveryMgr implements [tool.Bindings].
func (t *ToolBindings) GetDiscoveryMgr() *manager.DiscoveryManager {
	return t.discoveryMgr
}

func (t *ToolBindings) FindToolsByKeywords(keywords []string) ([]tool.Metadata, error) {
	return t.toolMgr.FindToolsByKeywords(keywords)
}

func (t *ToolBindings) RefreshCronJob(ctx *context.ExecuteContext) {
	t.eventBus.Publish(ctx, event_bus.RefreshCronJobEvent, &event_bus.RefreshCronJobEventPayload{})
}

func (t *ToolBindings) CreateCronJob(ctx *context.ExecuteContext, sessionID, name, prompt, cronExpr string) (string, error) {
	if cronExpr == "" {
		return "", errors.New("cron expression cannot be empty")
	}

	newJob, err := t.cronJobMgr.Create(ctx, name, sessionID, prompt, cronExpr)
	if err != nil {
		return "", err
	}

	return newJob.ID.String(), nil
}

func (t *ToolBindings) GetCronJobManager() *cronjob_manager.CronJobManager {
	return t.cronJobMgr
}

func (t *ToolBindings) CallTool(ctx *context.ExecuteContext, toolName string, args string) *tool.Result {
	return t.toolMgr.Call(ctx, &tool.Call{
		CallID: uuid.New().String(),
		Name:   toolName,
		Args:   args,
	})
}

func (t *ToolBindings) FindToolsByCategories(_ *context.ExecuteContext, categories []string) ([]tool.Metadata, error) {
	list := make([]tool.Metadata, 0)
	for _, item := range registry.Registry {

		categoriesMatch := false
		for _, category := range item.Tool.GetCategories() {
			if slices.Contains(categories, category) {
				categoriesMatch = true
				break
			}
		}

		if categoriesMatch {
			list = append(list, tool.Metadata{
				Name:        item.Tool.GetName(),
				Description: item.Tool.GetDescription(),
				Args:        item.Tool.GetArgsSchema(),
			})
		}
	}
	return list, nil
}

func (t *ToolBindings) RequestToolApproval(ctx *context.ExecuteContext, title string) bool {
	if ctx == nil || ctx.Context == nil || ctx.MessageChannel == nil {
		t.logger.Error("Cannot request tool approval without an active session channel", "title", title)
		return false
	}

	request := t.approvalMgr.Request(ctx.SessionID, title)
	select {
	case ctx.MessageChannel <- request:
	case <-ctx.Done():
		t.approvalMgr.Cancel(request.RequestID)
		t.logger.Error("Approval request aborted before reaching the user", "sessionID", ctx.SessionID, "title", title)
		return false
	}

	accepted, err := t.approvalMgr.Wait(ctx, request.RequestID)
	if err != nil {
		t.logger.Error("Error when request user approve action", "sessionID", ctx.SessionID, "title", title, "err", err)
		return false
	}

	return accepted
}
