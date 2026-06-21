package chat

import (
	tea "charm.land/bubbletea/v2"
	"github.com/smtdfc/nagare/core/messages"
)

type StreamMsg struct {
	Role      string
	ChunkType string
	Content   string
}
type StreamDoneMsg struct{}

func waitForMessage(ch <-chan messages.Message) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return StreamDoneMsg{}
		}

		switch chunk := msg.(type) {
		case *messages.StreamErrorMessage:
			return StreamMsg{
				Role:      "system",
				ChunkType: "error",
				Content:   chunk.Cause,
			}
		case *messages.ToolCallMessage:
			return StreamMsg{
				Role:      "agent",
				ChunkType: "tool_call",
				Content:   chunk.FunctionName,
			}
		case *messages.ToolResultMessage:
			return StreamMsg{
				Role:      "agent",
				ChunkType: "tool_call_result",
				Content:   chunk.CallID,
			}
		case *messages.TextMessage:
			return StreamMsg{
				Role:      "agent",
				ChunkType: "text",
				Content:   chunk.Content,
			}
		case *messages.ReasoningMessage:
			return StreamMsg{
				Role:      "agent",
				ChunkType: "reasoning",
				Content:   chunk.Content,
			}
		case *messages.AgentResponseDoneMessage:
			return StreamMsg{
				Role:      "system",
				ChunkType: "done",
				Content:   "",
			}
		}

		return StreamMsg{Role: "agent", ChunkType: "unknown"}
	}
}
