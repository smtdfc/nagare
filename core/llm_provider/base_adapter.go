package llm_provider

import (
	"context"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/message"
)

type LLMProviderAdapter interface {
	GetModels(context.Context) ([]string, error)
	Send(ctx context.Context, model string, inputs message.ListMessage, tools tool.ListTool) (message.MessageReadOnlyChannel, error)
}
