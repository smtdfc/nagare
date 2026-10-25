package hil

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/pkgs/messages"
)

const DefaultApprovalTimeout = 5 * time.Minute

var (
	ErrApprovalNotFound        = errors.New("approval request not found")
	ErrApprovalSessionMismatch = errors.New("approval request does not belong to this session")
	ErrApprovalTimedOut        = errors.New("approval request timed out")
	ErrApprovalRejected        = errors.New("tool execution was rejected by the user")
)

type pendingApproval struct {
	sessionID string
	response  chan bool
	resolved  bool
}

// Manager owns in-flight approval requests. A request is single-use and is
// always removed after a decision, cancellation, or timeout.
type Manager struct {
	mu      sync.Mutex
	pending map[string]*pendingApproval
	timeout time.Duration
}

func (m *Manager) Request(sessionID, title string) *messages.ApprovalRequestMessage {
	requestID := uuid.NewString()
	m.mu.Lock()
	m.pending[requestID] = &pendingApproval{
		sessionID: sessionID,
		response:  make(chan bool, 1),
	}
	m.mu.Unlock()

	return messages.NewApprovalRequestMessage(sessionID, requestID, title)
}

func (m *Manager) Wait(ctx context.Context, requestID string) (bool, error) {
	m.mu.Lock()
	pending, ok := m.pending[requestID]
	m.mu.Unlock()
	if !ok {
		return false, ErrApprovalNotFound
	}

	timer := time.NewTimer(m.timeout)
	defer timer.Stop()
	defer m.remove(requestID)

	select {
	case accepted := <-pending.response:
		return accepted, nil
	case <-ctx.Done():
		return false, ctx.Err()
	case <-timer.C:
		return false, ErrApprovalTimedOut
	}
}

func (m *Manager) Respond(requestID, sessionID string, accept bool) error {
	m.mu.Lock()
	pending, ok := m.pending[requestID]
	if !ok {
		m.mu.Unlock()
		return ErrApprovalNotFound
	}
	if pending.sessionID != sessionID {
		m.mu.Unlock()
		return ErrApprovalSessionMismatch
	}
	if pending.resolved {
		m.mu.Unlock()
		return ErrApprovalNotFound
	}
	pending.resolved = true
	m.mu.Unlock()

	pending.response <- accept
	return nil
}

func (m *Manager) remove(requestID string) {
	m.mu.Lock()
	delete(m.pending, requestID)
	m.mu.Unlock()
}

// Cancel discards a pending request that never reached the user, so it cannot
// leak in the pending map when delivery is aborted.
func (m *Manager) Cancel(requestID string) {
	m.remove(requestID)
}

// @Injectable
func NewManager() *Manager {
	return &Manager{
		pending: make(map[string]*pendingApproval),
		timeout: DefaultApprovalTimeout,
	}
}
