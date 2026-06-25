package agent

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/messages"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
)

const NAGARE_LIST_MESSAGE_SIZE_LIMIT = 10

var DEFAULT_MESSAGES = messages.ListMessage{
	SYSTEM_PROMPT,
	DEVELOPER_PROMPT,
}

type SessionManager struct {
	sync.RWMutex
	data   map[string]messages.ListMessage
	logger *slog.Logger
}

func (s *SessionManager) GetHistory(sessionID string, limit int) messages.ListMessage {
	s.RLock()
	defer s.RUnlock()

	history, exists := s.data[sessionID]
	if !exists {
		return DEFAULT_MESSAGES
	}

	return history
}

func (s *SessionManager) SaveHistory(sessionID string, history messages.ListMessage) {
	s.logger.Info(fmt.Sprintf("History has been updated for session %s", sessionID))
	s.Lock()
	defer s.Unlock()
	s.data[sessionID] = history
}

func (s *SessionManager) CreateSessionID() string {
	id := uuid.New()
	return id.String()
}

func NewSessionManager() *SessionManager {
	return &SessionManager{data: make(map[string]messages.ListMessage), logger: nagare_logger.GetLogger("Session Manager")}
}
