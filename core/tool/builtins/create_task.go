package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type CreateTaskInput struct {
	Name       string `json:"name"`
	Prompt     string `json:"prompt"`
	Repeat     bool   `json:"repeat"`
	RepeatRule string `json:"repeat_rule"` // includes: "no_repeat", "daily"
	TriggerBy  string `json:"trigger_by"`  // includes: "schedule"
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type CreateTaskOutput struct {
	TaskID string `json:"task_id"`
}

var CreateTaskTool = tool.DefineTool(
	"create_task_tool",
	"Create a scheduled task or an asynchronous background task assigned by the user. "+
		"Use this tool whenever the user wants to schedule a reminder, a delayed action, or an asynchronous job to be executed later. "+
		"Use 'name' for the task title, 'prompt' for the execution content/command, "+
		"'trigger_by' set to 'scheduled', 'repeat' as true/false, and 'repeat_rule' as 'no_repeat' or 'daily'. "+
		"If the task repeats indefinitely without an ending, leave 'end_time' empty or omit it. "+
		"For 'start_time' and 'end_time', format them strictly as 'YYYY-MM-DD HH:MM:SS' or RFC3339. "+
		"CRITICAL: If you need precise current time for calculating schedules, you MUST call 'time_tool' first to fetch it before setting the times.",
	func(ctx *context.ExecuteContext, args *CreateTaskInput, bindings tool.Bindings) (*CreateTaskOutput, error) {
		taskID, err := bindings.CreateTask(
			ctx,
			ctx.SessionID,
			args.Name,
			args.Prompt,
			args.TriggerBy,
			args.Repeat,
			args.RepeatRule,
			args.StartTime,
			args.EndTime,
		)

		if err != nil {
			return nil, err
		}

		bindings.RefreshTask(ctx)
		return &CreateTaskOutput{
			TaskID: taskID,
		}, nil
	},
)
