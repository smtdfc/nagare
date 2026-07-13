package llm

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/shared/messages"
)

type ChatCallback func(messages.Message)
type MessageChannel chan messages.Message
type LLMProviderAdapter interface {
	Chat(string, *context.ExecuteContext, messages.ListMessage) (MessageChannel, error)
}
