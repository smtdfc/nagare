package domains

type Tool interface {
	GetName() string
	GetArgs() string
	GetDescription() string
	Execute(ctx Context, args string) (string, error)
}

type ListTool []Tool

type ToolCallStatus string

const (
	TOOL_CALL_SUCCESS ToolCallStatus = "TOOL_CALL_SUCCESS"
	TOOL_CALL_FAILED  ToolCallStatus = "TOOL_CALL_FAILED"
	TOOL_CALL_PENDING ToolCallStatus = "TOOL_CALL_PENDING"
)

type ToolResult struct {
	Status ToolCallStatus
	Result string
	Error  string
	CallID string
}

func (r *ToolResult) Success(result string) *ToolResult {
	r.Status = TOOL_CALL_SUCCESS
	r.Result = result

	return r
}

func (r *ToolResult) Failed(err string) *ToolResult {
	r.Status = TOOL_CALL_FAILED
	r.Error = err

	return r
}

func NewToolResult(callID string) *ToolResult {
	return &ToolResult{
		Result: "",
		Error:  "",
		Status: TOOL_CALL_PENDING,
		CallID: callID,
	}
}

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
