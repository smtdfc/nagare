package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/smtdfc/nagare/shared/messages"
)

func ParseAgentMessage(raw json.RawMessage) (messages.Message, error) {
	var peek struct {
		Type messages.MessageType `json:"type"`
	}
	if err := json.Unmarshal(raw, &peek); err != nil {
		return nil, err
	}

	switch peek.Type {
	case messages.AGENT_RESPONSE:
		var msg messages.AgentResponse
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.REASONING_MESSAGE:
		var msg messages.Reasoning
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.RESPONSE_STARTED_MESSAGE:
		var msg messages.ResponseStarted
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.RESPONSE_COMPLETED_MESSAGE:
		var msg messages.ResponseCompleted
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.RESPONSE_FAILED_MESSAGE:
		var msg messages.ResponseFailed
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.TEXT_MESSAGE:
		var msg messages.Text
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.TOOL_CALL_MESSAGE:
		var msg messages.ToolCall
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	case messages.TOOL_RESULT_MESSAGE:
		var msg messages.ToolCallResult
		if err := json.Unmarshal(raw, &msg); err != nil {
			return nil, err
		}
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown message type: %s", peek.Type)
	}
}
