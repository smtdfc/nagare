package messages

type MessageType string

func (m MessageType) ToString() string {
	return string(m)
}

const (
	ResponseStartedMessageType   MessageType = "RESPONSE_STARTED_MESSAGE"
	ResponseCompletedMessageType MessageType = "RESPONSE_COMPLETED_MESSAGE"
	ResponseFailedMessageType    MessageType = "RESPONSE_FAILED_MESSAGE"
	TextMessageType              MessageType = "TEXT_MESSAGE"
	ToolCallMessageType          MessageType = "TOOL_CALL_MESSAGE"
	ToolResultMessageType        MessageType = "TOOL_RESULT_MESSAGE"
	ReasoningMessageType         MessageType = "REASONING_MESSAGE"
	AgentStartedMessageType      MessageType = "AGENT_STARTED_MESSAGE"
	AgentCompletedMessageType    MessageType = "AGENT_COMPLETED_MESSAGE"
	AgentErrorMessageType        MessageType = "AGENT_ERROR_MESSAGE"
)
