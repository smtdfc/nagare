package tool

type ToolCall struct {
	Name   string
	Args   string
	CallID string
}

func NewToolCall(name, args, callID string) *ToolCall {
	return &ToolCall{
		Name:   name,
		Args:   args,
		CallID: callID,
	}
}

type ListToolCall []*ToolCall

var EMPTY_LIST_TOOL_CALL = ListToolCall{}
