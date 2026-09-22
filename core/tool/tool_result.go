package tool

import "github.com/smtdfc/nagare/shared/messages"

type Result struct {
	callID    string
	name      string
	IsSuccess bool
	Result    string
	err       error
}

func (t *Result) ToMessage() *messages.ToolResultMessage {
	Result := t.Result
	if !t.IsSuccess {
		Result = t.err.Error()
	}

	return messages.NewToolResultMessage(
		t.callID,
		t.name,
		Result,
	)
}

type ResultBuilder struct {
	Result *Result
}

func NewToolResultBuilder(callID, name string) *ResultBuilder {
	return &ResultBuilder{
		Result: &Result{
			callID:    callID,
			name:      name,
			IsSuccess: true,
		},
	}
}

func (b *ResultBuilder) Success(Result string) *ResultBuilder {
	b.Result.IsSuccess = true
	b.Result.Result = Result
	b.Result.err = nil
	return b
}

func (b *ResultBuilder) Failure(err error) *ResultBuilder {
	b.Result.IsSuccess = false
	b.Result.err = err
	if err != nil {
		b.Result.Result = err.Error()
	}
	return b
}

func (b *ResultBuilder) SetResult(Result string) *ResultBuilder {
	b.Result.Result = Result
	return b
}

func (b *ResultBuilder) Build() *Result {
	return b.Result
}
