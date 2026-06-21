package agent

import (
	"fmt"
	"sync"

	"github.com/smtdfc/nagare/core/messages"
)

type SessionManager struct {
	sync.RWMutex
	data map[string]messages.ListMessage
}

func (s *SessionManager) GetHistory(sessionID string) messages.ListMessage {
	s.RLock()
	defer s.RUnlock()
	if history, exists := s.data[sessionID]; exists {
		appLogger.Info(fmt.Sprintf("Successfully retrieved chat history for session %s", sessionID))
		return history
	}
	return messages.ListMessage{SYSTEM_PROMPT}
}

func (s *SessionManager) SaveHistory(sessionID string, history messages.ListMessage) {
	appLogger.Info(fmt.Sprintf("History has been updated for session %s", sessionID))
	s.Lock()
	defer s.Unlock()
	s.data[sessionID] = history
}

func NewSessionManager() *SessionManager {
	return &SessionManager{data: make(map[string]messages.ListMessage)}
}
