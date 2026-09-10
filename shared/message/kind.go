package message

type Kind string

func (m Kind) ToString() string {
	return string(m)
}

const (
	ResponseStartedMessageKind   Kind = "RESPONSE_STARTED_MESSAGE"
	ResponseCompletedMessageKind Kind = "RESPONSE_COMPLETED_MESSAGE"
	ResponseFailedMessageKind    Kind = "RESPONSE_FAILED_MESSAGE"
	TextMessageKind              Kind = "TEXT_MESSAGE"
	ToolCallMessageKind          Kind = "TOOL_CALL_MESSAGE"
	ToolResultMessageKind        Kind = "TOOL_RESULT_MESSAGE"
	ReasoningMessageKind         Kind = "REASONING_MESSAGE"
	AgentStartedMessageKind      Kind = "AGENT_STARTED_MESSAGE"
	AgentCompletedMessageKind    Kind = "AGENT_COMPLETED_MESSAGE"
	AgentErrorMessageKind        Kind = "AGENT_ERROR_MESSAGE"
)
