package main

import (
	"context"
	"testing"
	"time"

	"github.com/mymmrac/telego"
)

func TestTelegramQueueFlow(t *testing.T) {
	tp := &TelegramPlugin{
		channels:     make(map[string]struct{}),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
	}

	tp.setChannel("123456")
	tp.clearChannel("123456")
	if tp.hasChannel("123456") {
		t.Fatalf("expected channel not to exist")
	}

	// Test processQueue with empty queue terminates loop immediately
	tp.processQueue(context.Background(), "123456")
	tp.mu.Lock()
	processing := tp.isProcessing["123456"]
	tp.mu.Unlock()
	if processing {
		t.Fatalf("expected processing=false for empty queue")
	}
}
