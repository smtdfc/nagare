package messages

import (
	"testing"
)

func TestMessageTypeToString(t *testing.T) {
	mt := TextMessageType
	if mt.ToString() != "TEXT_MESSAGE" {
		t.Fatalf("expected string 'TEXT_MESSAGE', got %q", mt.ToString())
	}
}

func TestMessageTypeToString_CompactTypes(t *testing.T) {
	tests := []struct {
		mt       MessageType
		expected string
	}{
		{ConversationSummaryMessageType, "CONVERSATION_SUMMARY_MESSAGE"},
		{CompactSuccessMessageType, "COMPACT_SUCCESS_MESSAGE"},
	}

	for _, tc := range tests {
		t.Run(string(tc.mt), func(t *testing.T) {
			if tc.mt.ToString() != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, tc.mt.ToString())
			}
		})
	}
}
