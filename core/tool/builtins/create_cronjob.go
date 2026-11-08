package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type CreateCronJobInput struct {
	Name              string `json:"name" jsonschema_description:"The name or title of the cronjob."`
	ActionDescription string `json:"action_description" jsonschema_description:"A self-contained task instruction to execute when this scheduled job triggers. This is NOT a request to create, modify, or schedule another cronjob. Execute the described task directly. Do not recreate this schedule."`
	CronExpr          string `json:"cron_expr" jsonschema_description:"Standard cron expression with 5 fields (minute hour day-of-month month day-of-week), e.g., '0 9 * * *' for every day at 9:00 AM, or '*/30 * * * *' for every 30 minutes, or standard macros like '@daily', '@hourly'."`
}

type CreateCronJobOutput struct {
	CronJobID string `json:"cronjob_id"`
}

var CreateCronJobTool = tool.DefineTool(
	"create_cronjob_tool",
	`Create a CRONJOB`,
	func(ctx *context.ExecuteContext, args *CreateCronJobInput, bindings tool.Bindings) (*CreateCronJobOutput, error) {
		jobID, err := bindings.CreateCronJob(
			ctx,
			ctx.SessionID,
			args.Name,
			args.ActionDescription,
			args.CronExpr,
		)

		if err != nil {
			return nil, err
		}

		bindings.RefreshCronJob(ctx)

		return &CreateCronJobOutput{
			CronJobID: jobID,
		}, nil
	},
	[]string{tool.CronJobManagementCategory},
)
