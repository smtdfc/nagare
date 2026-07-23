package messages

type ToolCall struct {
	Type   MessageType `json:"type"`
	Name   string      `json:"name"`
	Args   string      `json:"args"`
	CallID string      `json:"call_id"`
}

func NewToolCall(name, args, callID string) *ToolCall {
	return &ToolCall{
		Type:   TOOL_CALL_MESSAGE,
		Name:   name,
		Args:   args,
		CallID: callID,
	}
}

func (t *ToolCall) GetType() MessageType {
	return t.Type
}

type ToolCallResult struct {
	Type   MessageType `json:"type"`
	CallID string      `json:"call_id"`
	Result string      `json:"result"`
	Error  string      `json:"error"`
}

func NewToolCallResult(callID, result, err string) *ToolCallResult {
	return &ToolCallResult{
		Type:   TOOL_RESULT_MESSAGE,
		CallID: callID,
		Result: result,
		Error:  err,
	}
}

func (t *ToolCallResult) GetType() MessageType {
	return t.Type
}
