package llm_provider

import (
	"context"

	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/messages"
)

type LLMProviderAdapter interface {
	GetModels(context.Context) ([]string, error)
	Send(ctx context.Context, model string, inputs messages.ListMessage, tools tool.ListTool) (message.ReadOnlyChannel, error)
}
