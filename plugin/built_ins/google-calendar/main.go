package main

import (
	"context"
	_ "embed"

	"github.com/smtdfc/nagare/plugin/client"
)

var categories = []string{"google_calendar"}

//go:embed metadata.json
var metadata string

type FindTaskInput struct {
	Keyword string `json:"keyword"`
}

type FindTaskOutput struct {
	TaskName string `json:"task_name"`
}

var findTaskTool = client.DefineTool(
	"google_calendar:find_task",
	"Find a task in the Google Calendar.",
	func(ctx *context.Context, args *FindTaskInput) (*FindTaskOutput, error) {
		// Simulate finding a task (replace with actual logic)
		taskName := "Sample Task Name" // Replace with actual task name retrieved from Google Calendar

		return &FindTaskOutput{
			TaskName: taskName,
		}, nil
	},
	categories,
)

type CreateReminderInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // Format: YYYY-MM-DD
}

type CreateReminderOutput struct {
	ReminderID string `json:"reminder_id"`
}

var createReminderTool = client.DefineTool(
	"google_calendar:create_reminder",
	"Create a reminder in the Google Calendar.",
	func(ctx *context.Context, args *CreateReminderInput) (*CreateReminderOutput, error) {
		// Simulate creating a reminder (replace with actual logic)
		reminderID := "3432432565758745" // Replace with actual reminder ID returned from Google Calendar

		return &CreateReminderOutput{
			ReminderID: reminderID,
		}, nil
	},
	categories,
)

func OnStart(ctx context.Context, pluginClient *client.PluginClient) {
	pluginClient.RegisterToolCategories(ctx, categories)
	pluginClient.RegisterPluginTool(ctx, findTaskTool)
	pluginClient.RegisterPluginTool(ctx, createReminderTool)
}

func main() {
	ctx := context.Background()
	pluginClient := client.NewPlugin()
	_, err := pluginClient.LoadMetadata(metadata)
	if err != nil {
		pluginClient.Logger.Error("Load metadata error", "error", err)
		return
	}

	err = pluginClient.Start(
		ctx,
		func() {
			err := pluginClient.Handshake(ctx)
			if err != nil {
				pluginClient.Logger.Error("Handshake error", "error", err)
				return
			}

			OnStart(ctx, pluginClient)
		},
	)
	if err != nil {
		pluginClient.Logger.Error("Error", "error", err)
	}
}
