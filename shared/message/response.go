package message

import "github.com/smtdfc/nagare/shared/helpers"

type ResponseStartedMessage struct {
	ID string `json:"id"`
}

func (t *ResponseStartedMessage) GetKind() Kind {
	return ResponseStartedMessageKind
}

func NewResponseStartedMessage() *ResponseStartedMessage {
	return &ResponseStartedMessage{
		ID: helpers.GenerateUUID(),
	}
}

type ResponseCompletedMessage struct {
	ID string `json:"id"`
}

func (t *ResponseCompletedMessage) GetKind() Kind {
	return ResponseCompletedMessageKind
}

func NewResponseCompletedMessage() *ResponseCompletedMessage {
	return &ResponseCompletedMessage{
		ID: helpers.GenerateUUID(),
	}
}

type ResponseFailedMessage struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Cause string `json:"cause"`
}

func (t *ResponseFailedMessage) GetKind() Kind {
	return ResponseFailedMessageKind
}

func NewResponseFailedMessage(code string, cause string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:    helpers.GenerateUUID(),
		Cause: cause,
		Code:  code,
	}
}
