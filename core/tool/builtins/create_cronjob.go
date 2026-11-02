package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type CreateCronJobInput struct {
	Name     string `json:"name" jsonschema_description:"The name or title of the cronjob."`
	Prompt   string `json:"prompt" jsonschema_description:"The execution instructions or command that the agent will run when triggered."`
	CronExpr string `json:"cron_expr" jsonschema_description:"Standard cron expression with 5 fields (minute hour day-of-month month day-of-week), e.g., '0 9 * * *' for every day at 9:00 AM, or '*/30 * * * *' for every 30 minutes, or standard macros like '@daily', '@hourly'."`
}

type CreateCronJobOutput struct {
	CronJobID string `json:"cronjob_id"`
}

var CreateCronJobTool = tool.DefineTool(
	"create_cronjob_tool",
	"CRITICAL RESTRICTION: This tool is STRICTLY for scheduling recurring or automated actions, reminders, or future cronjob tasks using cron expressions. DO NOT use it as a substitute for general memory or note-taking. "+
		"Use this tool ONLY when the user explicitly requests to schedule something recurring or at a future time. "+
		"Use 'name' for the job title, 'prompt' for the execution instruction/content, and 'cron_expr' for the cron schedule expression. "+
		"CRITICAL: If you need precise current time for calculating schedules, you MUST call 'time_tool' first to fetch it before determining the expression.",
	func(ctx *context.ExecuteContext, args *CreateCronJobInput, bindings tool.Bindings) (*CreateCronJobOutput, error) {
		jobID, err := bindings.CreateCronJob(
			ctx,
			ctx.SessionID,
			args.Name,
			args.Prompt,
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
