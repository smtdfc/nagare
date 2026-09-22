package messages

import "github.com/smtdfc/nagare/shared/helpers"

type ResponseStartedMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (t *ResponseStartedMessage) GetMessageType() MessageType {
	return t.Type
}

func NewResponseStartedMessage() *ResponseStartedMessage {
	return &ResponseStartedMessage{
		ID:   helpers.GenerateUUID(),
		Type: ResponseStartedMessageType,
	}
}

type ResponseCompletedMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (t *ResponseCompletedMessage) GetMessageType() MessageType {
	return t.Type
}

func NewResponseCompletedMessage() *ResponseCompletedMessage {
	return &ResponseCompletedMessage{
		ID:   helpers.GenerateUUID(),
		Type: ResponseCompletedMessageType,
	}
}

type ResponseFailedMessage struct {
	ID    string      `json:"id"`
	Type  MessageType `json:"type"`
	Code  string      `json:"code"`
	Cause string      `json:"cause"`
}

func (t *ResponseFailedMessage) GetMessageType() MessageType {
	return t.Type
}

func NewResponseFailedMessage(code string, cause string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:    helpers.GenerateUUID(),
		Type:  ResponseFailedMessageType,
		Cause: cause,
		Code:  code,
	}
}
