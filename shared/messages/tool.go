package messages

type ToolCall struct {
	Name   string `json:"name"`
	Args   string `json:"args"`
	CallID string `json:"call_id"`
}

func (t *ToolCall) Kind() string {
	return "ToolCall"
}

type ToolCallResult struct {
	CallID string `json:"call_id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

func (t *ToolCallResult) Kind() string {
	return "ToolCallResult"
}
