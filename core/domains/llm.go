package domains

import (
	"github.com/smtdfc/nagare/shared/messages"
)

type ChatCallback func(messages.Message)
type MessageChannel chan messages.Message
type LLMProviderAdapter interface {
	Chat(string, Context, messages.ListMessage) (MessageChannel, error)
}
