package chat

import (
	"testing"
)

func TestChatWorker_Struct(t *testing.T) {
	// Verify ChatWorker struct fields
	worker := &Worker{}
	if worker.chatEventBus != nil {
		t.Errorf("expected nil initial chatEventBus")
	}
}
