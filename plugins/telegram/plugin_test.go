package main

import (
	"context"
	"testing"
)

func TestTelegramPluginInstance(t *testing.T) {
	tp := NewTelegramPlugin(nil, nil)
	if tp == nil {
		t.Fatalf("expected non-nil TelegramPlugin instance")
	}

	// Test sendTextMessage with empty text returns nil
	err := tp.sendTextMessage(context.Background(), 12345, "")
	if err != nil {
		t.Fatalf("expected nil error for empty message, got %v", err)
	}
}
