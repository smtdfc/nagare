package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/smtdfc/nagare/plugin/client"
)

var googleCalendarToolCategories = map[string]string{
	"google_calendar": "Google Calendar Tools",
}

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
	[]string{"google_calendar"},
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
		fmt.Println("called")
		return &CreateReminderOutput{
			ReminderID: reminderID,
		}, nil
	},
	[]string{"google_calendar"},
)

func OnStart(ctx context.Context, pluginClient *client.PluginClient) {
	// Register the tools with the plugin client
	err := pluginClient.RegisterPluginTool(ctx, findTaskTool)
	if err != nil {
		pluginClient.Logger.Error("Failed to register find task tool", "error", err)
		return
	}

	err = pluginClient.RegisterPluginTool(ctx, createReminderTool)
	if err != nil {
		pluginClient.Logger.Error("Failed to register create reminder tool", "error", err)
		return
	}

	// Register the tool categories with the plugin client
	err = pluginClient.RegisterToolCategories(ctx, googleCalendarToolCategories)
	if err != nil {
		pluginClient.Logger.Error("Failed to register tool categories", "error", err)
		return
	}
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
