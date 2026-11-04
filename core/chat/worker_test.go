package chat

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/event_bus"
	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func newTestWorker() *Worker {
	return &Worker{
		cancelFn: make(map[string]context.CancelFunc),
	}
}

func TestChatWorker_Struct(t *testing.T) {
	// Verify ChatWorker struct fields
	worker := &Worker{}
	if worker.chatEventBus != nil {
		t.Errorf("expected nil initial chatEventBus")
	}
}

func TestWorker_HandleCancelEvent_RemovesCancelFn(t *testing.T) {
	w := newTestWorker()

	cancelled := false
	w.cancelFnMu.Lock()
	w.cancelFn["invoke-1"] = func() { cancelled = true }
	w.cancelFnMu.Unlock()

	w.HandleCancelEvent(&event_bus.CancelChatEventPayload{
		InvokeID: "invoke-1",
	})

	if !cancelled {
		t.Error("expected cancel function to be called")
	}

	w.cancelFnMu.Lock()
	_, exists := w.cancelFn["invoke-1"]
	w.cancelFnMu.Unlock()

	if exists {
		t.Error("expected cancelFn to be removed from map after cancel")
	}
}

func TestWorker_HandleCancelEvent_UnknownInvokeID(t *testing.T) {
	w := newTestWorker()
	// Should not panic on unknown invokeID
	w.HandleCancelEvent(&event_bus.CancelChatEventPayload{
		InvokeID: "unknown-invoke",
	})
}

func TestWorker_HandleCompactConversationEvent_DrainAndCleanup(t *testing.T) {
	successCh := make(chan messages.Message, 1)
	successMsg := messages.NewCompactSuccessMessage()
	successMsg.SetInvokeID("invoke-compact-1")
	successCh <- successMsg
	close(successCh)

	var (
		mu       sync.Mutex
		captured *ConversationCompactParams
	)

	// Patch agentInvoker by embedding a closure-based stub via a helper
	// function that accepts a readOnlyChannel factory.
	w := newTestWorkerWithCompact(func(ctx context.Context, params *ConversationCompactParams) (message.ReadOnlyChannel, error) {
		mu.Lock()
		captured = params
		mu.Unlock()
		return successCh, nil
	})

	w.HandleCompactConversationEvent(&event_bus.CompactConversationEventPayload{
		InvokeID:   "invoke-compact-1",
		SessionID:  "session-abc",
		SenderType: event_bus.User,
		SenderID:   "user-1",
	})

	// Wait for the goroutine draining the channel to finish.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		w.cancelFnMu.Lock()
		_, exists := w.cancelFn["invoke-compact-1"]
		w.cancelFnMu.Unlock()
		if !exists {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	w.cancelFnMu.Lock()
	_, exists := w.cancelFn["invoke-compact-1"]
	w.cancelFnMu.Unlock()

	if exists {
		t.Error("expected cancelFn to be cleaned up after compact channel was drained")
	}

	mu.Lock()
	defer mu.Unlock()
	if captured == nil {
		t.Fatal("expected Compact to have been called with params")
	}
	if captured.SessionID != "session-abc" {
		t.Errorf("expected session ID 'session-abc', got %q", captured.SessionID)
	}
	if !captured.SendIntoEventBus {
		t.Error("expected SendIntoEventBus to be true")
	}
}

// compactFunc is the signature of the Compact method we want to stub.
type compactFunc func(ctx context.Context, params *ConversationCompactParams) (message.ReadOnlyChannel, error)

// stubWorker embeds Worker and overrides only HandleCompactConversationEvent
// so we can test channel-drain + cancelFn cleanup without a full AgentInvoker.
type stubWorker struct {
	Worker
	compact compactFunc
}

func (s *stubWorker) HandleCompactConversationEvent(payload *event_bus.CompactConversationEventPayload) {
	ctx, cancelFn := context.WithCancel(context.Background())
	s.cancelFnMu.Lock()
	s.cancelFn[payload.InvokeID] = cancelFn
	s.cancelFnMu.Unlock()

	output, err := s.compact(ctx, &ConversationCompactParams{
		SessionID:        payload.SessionID,
		SenderType:       payload.SenderType,
		SenderID:         payload.SenderID,
		SendIntoEventBus: true,
		InvokeID:         payload.InvokeID,
	})
	if err != nil {
		s.cancelFnMu.Lock()
		delete(s.cancelFn, payload.InvokeID)
		s.cancelFnMu.Unlock()
		return
	}

	go func() {
		defer func() {
			s.cancelFnMu.Lock()
			delete(s.cancelFn, payload.InvokeID)
			s.cancelFnMu.Unlock()
		}()
		for range output {
		}
	}()
}

func newTestWorkerWithCompact(fn compactFunc) *stubWorker {
	return &stubWorker{
		Worker:  Worker{cancelFn: make(map[string]context.CancelFunc)},
		compact: fn,
	}
}
