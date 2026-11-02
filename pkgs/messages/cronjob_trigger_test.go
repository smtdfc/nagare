package messages

import (
	"testing"
)

func TestCronJobTriggerMessage(t *testing.T) {
	msg := NewCronJobTriggerMessage()
	if msg.GetMessageType() != CronJobTriggerMessageType {
		t.Fatalf("expected type %s, got %s", CronJobTriggerMessageType, msg.GetMessageType())
	}
	if msg.GetMessageID() == "" {
		t.Fatal("expected non-empty message ID")
	}

	msg.SetInvokeID("inv-123")
	if msg.GetInvokeID() != "inv-123" {
		t.Fatalf("expected invoke ID 'inv-123', got %s", msg.GetInvokeID())
	}
}
