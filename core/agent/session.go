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
		return messages.ListMessage{SYSTEM_PROMPT}
	}

	// apply limit
	if limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}

	// pre-allocate: 1 system + history size
	n := len(history) + 1
	result := make(messages.ListMessage, 0, n)

	result = append(result, SYSTEM_PROMPT)
	result = append(result, history...)

	return result
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
