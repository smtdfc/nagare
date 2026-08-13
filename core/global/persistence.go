package global

import "github.com/smtdfc/nagare/core/persistence/database/repositories"

var GlobalLlmRepository *repositories.LLMProviderRepository
var GlobalKVRepository *repositories.KVRepository
var GlobalSessionRepository *repositories.SessionRepository
var GlobalMessageRepository *repositories.MessageRepository
