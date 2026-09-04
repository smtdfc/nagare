package tool

import "github.com/smtdfc/nagare/shared/message"

type ToolResult struct {
	callID    string
	name      string
	IsSuccess bool
	Result    string
	err       error
}

func (t *ToolResult) ToMessage() *message.ToolResultMessage {
	Result := t.Result
	if !t.IsSuccess {
		Result = t.err.Error()
	}

	return message.NewToolResultMessage(
		t.callID,
		t.name,
		Result,
	)
}

type ToolResultBuilder struct {
	Result *ToolResult
}

func NewToolResultBuilder(callID, name string) *ToolResultBuilder {
	return &ToolResultBuilder{
		Result: &ToolResult{
			callID:    callID,
			name:      name,
			IsSuccess: true,
		},
	}
}

func (b *ToolResultBuilder) Success(Result string) *ToolResultBuilder {
	b.Result.IsSuccess = true
	b.Result.Result = Result
	b.Result.err = nil
	return b
}

func (b *ToolResultBuilder) Failure(err error) *ToolResultBuilder {
	b.Result.IsSuccess = false
	b.Result.err = err
	if err != nil {
		b.Result.Result = err.Error()
	}
	return b
}

func (b *ToolResultBuilder) SetResult(Result string) *ToolResultBuilder {
	b.Result.Result = Result
	return b
}

func (b *ToolResultBuilder) Build() *ToolResult {
	return b.Result
}
