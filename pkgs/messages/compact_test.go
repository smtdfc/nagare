package messages

import (
	"testing"
)

func TestConversationSummaryMessage(t *testing.T) {
	const summaryText = "The user asked about X, the agent explained Y."
	msg := NewConversationSummaryMessage(summaryText)
	if msg == nil {
		t.Fatal("expected non-nil ConversationSummaryMessage")
	}

	if msg.GetMessageID() == "" {
		t.Fatal("expected non-empty generated UUID for message ID")
	}

	if msg.GetMessageType() != ConversationSummaryMessageType {
		t.Fatalf("expected ConversationSummaryMessageType, got %v", msg.GetMessageType())
	}

	if msg.Summary != summaryText {
		t.Fatalf("expected summary %q, got %q", summaryText, msg.Summary)
	}

	// InvokeID should be empty before being set
	if msg.GetInvokeID() != "" {
		t.Fatalf("expected empty InvokeID before set, got %q", msg.GetInvokeID())
	}

	msg.SetInvokeID("invoke-summary-001")
	if msg.GetInvokeID() != "invoke-summary-001" {
		t.Fatalf("expected invoke ID 'invoke-summary-001', got %q", msg.GetInvokeID())
	}
}

func TestConversationSummaryMessage_EmptySummary(t *testing.T) {
	msg := NewConversationSummaryMessage("")
	if msg == nil {
		t.Fatal("expected non-nil ConversationSummaryMessage with empty summary")
	}
	if msg.Summary != "" {
		t.Fatalf("expected empty summary, got %q", msg.Summary)
	}
}

func TestConversationSummaryMessage_UniqueIDs(t *testing.T) {
	msg1 := NewConversationSummaryMessage("summary 1")
	msg2 := NewConversationSummaryMessage("summary 2")
	if msg1.GetMessageID() == msg2.GetMessageID() {
		t.Fatal("expected unique IDs for distinct ConversationSummaryMessage instances")
	}
}

func TestCompactSuccessMessage(t *testing.T) {
	msg := NewCompactSuccessMessage()
	if msg == nil {
		t.Fatal("expected non-nil CompactSuccessMessage")
	}

	if msg.GetMessageID() == "" {
		t.Fatal("expected non-empty generated UUID for message ID")
	}

	if msg.GetMessageType() != CompactSuccessMessageType {
		t.Fatalf("expected CompactSuccessMessageType, got %v", msg.GetMessageType())
	}

	// InvokeID should be empty before being set
	if msg.GetInvokeID() != "" {
		t.Fatalf("expected empty InvokeID before set, got %q", msg.GetInvokeID())
	}

	msg.SetInvokeID("invoke-compact-001")
	if msg.GetInvokeID() != "invoke-compact-001" {
		t.Fatalf("expected invoke ID 'invoke-compact-001', got %q", msg.GetInvokeID())
	}
}

func TestCompactSuccessMessage_UniqueIDs(t *testing.T) {
	msg1 := NewCompactSuccessMessage()
	msg2 := NewCompactSuccessMessage()
	if msg1.GetMessageID() == msg2.GetMessageID() {
		t.Fatal("expected unique IDs for distinct CompactSuccessMessage instances")
	}
}
