package message

import "github.com/smtdfc/nagare/shared/helpers"

type ResponseStartedMessage struct {
	ID string `json:"id"`
}

func (t *ResponseStartedMessage) GetKind() MessageKind {
	return RESPONSE_STARTED_MESSAGE
}

func NewResponseStartedMessage() *ResponseStartedMessage {
	return &ResponseStartedMessage{
		ID: helpers.GenerateUUID(),
	}
}

type ResponseCompletedMessage struct {
	ID string `json:"id"`
}

func (t *ResponseCompletedMessage) GetKind() MessageKind {
	return RESPONSE_COMPLETED_MESSAGE
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

func (t *ResponseFailedMessage) GetKind() MessageKind {
	return RESPONSE_FAILED_MESSAGE
}

func NewResponseFailedMessage(code string, cause string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:    helpers.GenerateUUID(),
		Cause: cause,
		Code:  code,
	}
}
