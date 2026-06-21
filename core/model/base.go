package model

import (
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/messages"
)

var appLogger = logger.GetLogger()

type ChatModelConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

type MessageCallback func(messages.Message)
type ChatModel interface {
	Chat(ctx context.ExecuteContext, messages messages.ListMessage, cb MessageCallback) error
}
