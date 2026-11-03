package main

import (
	"testing"
	"time"

	"github.com/mymmrac/telego"
)

func TestTelegramChannelManagement(t *testing.T) {
	tp := &TelegramPlugin{
		channels:     make(map[string]struct{}),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
		approvals:    make(map[string]pendingToolApproval),
	}

	if tp.hasChannel("123456") {
		t.Fatalf("expected channel not to exist")
	}

	tp.setChannel("123456")
	if !tp.hasChannel("123456") {
		t.Fatalf("expected channel to exist after setChannel")
	}

	tp.clearChannel("123456")
	if tp.hasChannel("123456") {
		t.Fatalf("expected channel not to exist after clearChannel")
	}

	tp.setChannel("123456")
	tp.setPendingApproval("req-1", "123456", 42)
	if _, ok := tp.getPendingApproval("req-1", "123456"); !ok {
		t.Fatalf("expected pending approval for channel")
	}
	if _, ok := tp.getPendingApproval("req-1", "other"); ok {
		t.Fatalf("expected pending approval to be scoped to its channel")
	}

	tp.releaseChannel("123456")
	if tp.hasChannel("123456") {
		t.Fatalf("expected channel to be released")
	}
	if _, ok := tp.getPendingApproval("req-1", "123456"); ok {
		t.Fatalf("expected pending approvals to be cleared on release")
	}

	// Test finishProcessing
	tp.isProcessing["123456"] = true
	tp.finishProcessing("123456")
	tp.mu.Lock()
	processing := tp.isProcessing["123456"]
	tp.mu.Unlock()
	if processing {
		t.Fatalf("expected isProcessing=false after finishProcessing")
	}
}
