package main

import (
	"testing"
	"time"

	"github.com/mymmrac/telego"
)

func TestTelegramHandlerChunkParsing(t *testing.T) {
	tp := &TelegramPlugin{
		channels:     make(map[string]struct{}),
		idleTimers:   make(map[string]*time.Timer),
		queues:       make(map[string][]telego.Update),
		isProcessing: make(map[string]bool),
	}

	// Test malformed chunk does not panic
	tp.OnReceivedChatMessage("12345", `{invalid_chunk`)

	// Test TextMessageType chunk accumulation
	validChunk := `{"type":"TEXT_MESSAGE","content":"Hello world"}`
	tp.OnReceivedChatMessage("12345", validChunk)
	buf := tp.getOrCreateBuffer("12345")
	if buf.String() != "Hello world" {
		t.Fatalf("expected buffered content 'Hello world', got %q", buf.String())
	}
}
