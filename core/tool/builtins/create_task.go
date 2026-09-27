package declarations

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type CreateTaskInput struct {
	Name       string `json:"name"`
	Prompt     string `json:"prompt"`
	Repeat     bool   `json:"repeat"`
	RepeatRule string `json:"repeat_rule" jsonschema_description:"Defines the repetition rule for the task. Options include 'no_repeat' for a one-time task and 'daily' for a task that repeats every day."`
	TriggerBy  string `json:"trigger_by" jsonschema_description:"Specifies the trigger mechanism for the task. Currently, only 'scheduled' is supported, indicating that the task will be executed based on a defined schedule."`
	StartTime  string `json:"start_time" jsonschema_description:"Specifies the start time for the task. This is relevant for tasks that have a defined start point. The format should be 'YYYY-MM-DD HH:MM:SS' or RFC3339. If the task is meant to run immediately or without a specific start time, this field can be left empty."`
	EndTime    string `json:"end_time" jsonschema_description:"Specifies the end time for the task. This is relevant for tasks that have a defined duration or end point. The format should be 'YYYY-MM-DD HH:MM:SS' or RFC3339. If the task is meant to run indefinitely or until manually stopped, this field can be left empty."`
}

type CreateTaskOutput struct {
	TaskID string `json:"task_id"`
}

var CreateTaskTool = tool.DefineTool(
	"create_task_tool",
	"CRITICAL RESTRICTION: This tool is STRICTLY for scheduling automated actions, reminders, or future cronjob tasks. DO NOT use it as a substitute for general memory, note-taking, or storing arbitrary information. "+
		"Use this tool ONLY when the user explicitly requests to schedule something for later. "+
		"Use 'name' for the task title, 'prompt' for the execution content/command, "+
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
	[]string{tool.TaskManagementCategory},
)
