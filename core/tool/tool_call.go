package tool

type ToolCall struct {
	CallID string
	Name   string
	Args   string
}

type ListToolCall []*ToolCall

func NewToolCall(callID, name, args string) *ToolCall {
	return &ToolCall{
		CallID: callID,
		Name:   name,
		Args:   args,
	}
}
