package event_bus

import (
	"testing"
)

func TestEventPayload_GetEventType(t *testing.T) {
	// Verify SenderType constants
	if User != "user" || Plugin != "plugin" || System != "system" {
		t.Errorf("unexpected SenderType constants")
	}

	// Verify EventType values
	if SendChatMessageEvent != "event:chat:send" || ChatMessageChunkEvent != "event:chat:chunk" || RefreshTaskEvent != "event:task:refresh" {
		t.Errorf("unexpected EventType constants")
	}

	// Verify GetEventType for each payload struct
	sendPayload := &SendChatMessageEventPayload{RequestID: "req-1"}
	if sendPayload.GetEventType() != SendChatMessageEvent {
		t.Errorf("expected SendEvent, got %s", sendPayload.GetEventType())
	}

	chunkPayload := &ChatMessageChunkEventPayload{RequestID: "req-2"}
	if chunkPayload.GetEventType() != ChatMessageChunkEvent {
		t.Errorf("expected ChunkEvent, got %s", chunkPayload.GetEventType())
	}

	refreshPayload := &RefreshTaskEventPayload{}
	if refreshPayload.GetEventType() != RefreshTaskEvent {
		t.Errorf("expected RefreshTaskEvent, got %s", refreshPayload.GetEventType())
	}

	toolCallPayload := &PluginToolCallEventPayload{RequestID: "req-3"}
	if toolCallPayload.GetEventType() != PluginToolCallEvent {
		t.Errorf("expected PluginToolCallEvent, got %s", toolCallPayload.GetEventType())
	}

	toolResultPayload := &PluginToolCallResultEventPayload{RequestID: "req-4"}
	if toolResultPayload.GetEventType() != PluginToolCallResultEvent {
		t.Errorf("expected PluginToolCallResultEvent, got %s", toolResultPayload.GetEventType())
	}
}
