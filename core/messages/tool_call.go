package messages

type ToolCallMessage struct {
	CallID       string `json:"call_id,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Args         string `json:"args,omitempty"`
}

func (m *ToolCallMessage) GetType() string {
	return "ToolCall"
}

type ToolResultMessage struct {
	CallID string
	Result string
	Error  error
}

func (m *ToolResultMessage) GetType() string {
	return "ToolCallResult"
}
