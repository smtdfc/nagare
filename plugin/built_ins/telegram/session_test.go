package main

import (
	"testing"
	"time"

	"github.com/mymmrac/telego"
)

func TestTelegramSessionManagement(t *testing.T) {
	tp := &TelegramPlugin{
		channels:     make(map[string]string),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
	}

	// Test set and get sessionID
	tp.setSessionID("123456", "session-abc")
	sessID, exists := tp.getSessionID("123456")
	if !exists || sessID != "session-abc" {
		t.Fatalf("expected session-abc, got %q (exists=%v)", sessID, exists)
	}

	// Test getChatIDBySessionID
	strChatID, intChatID, ok := tp.getChatIDBySessionID("session-abc")
	if !ok || strChatID != "123456" || intChatID != 123456 {
		t.Fatalf("chat id lookup mismatch: str=%s, int=%d, ok=%v", strChatID, intChatID, ok)
	}

	// Test non-existent session
	_, _, ok = tp.getChatIDBySessionID("non-existent")
	if ok {
		t.Fatalf("expected ok=false for non-existent session")
	}

	// Test finishProcessing
	tp.isProcessing["123456"] = true
	tp.finishProcessing("123456", "session-abc")
	tp.mu.Lock()
	processing := tp.isProcessing["123456"]
	tp.mu.Unlock()
	if processing {
		t.Fatalf("expected isProcessing=false after finishProcessing")
	}
}
