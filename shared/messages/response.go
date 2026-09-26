package messages

import "github.com/smtdfc/nagare/shared/helpers"

type ResponseStartedMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (m *ResponseStartedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseStartedMessage) GetMessageType() MessageType {
	return m.Type
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

func (m *ResponseCompletedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseCompletedMessage) GetMessageType() MessageType {
	return m.Type
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

func (m *ResponseFailedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseFailedMessage) GetMessageType() MessageType {
	return m.Type
}

func NewResponseFailedMessage(code string, cause string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:    helpers.GenerateUUID(),
		Type:  ResponseFailedMessageType,
		Cause: cause,
		Code:  code,
	}
}
