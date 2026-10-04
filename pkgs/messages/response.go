package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type ResponseStartedMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

func (m *ResponseStartedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseStartedMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ResponseStartedMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ResponseStartedMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewResponseStartedMessage() *ResponseStartedMessage {
	return &ResponseStartedMessage{
		ID:   helpers.GenerateUUID(),
		Type: ResponseStartedMessageType,
	}
}

type ResponseCompletedMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

func (m *ResponseCompletedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseCompletedMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ResponseCompletedMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ResponseCompletedMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewResponseCompletedMessage() *ResponseCompletedMessage {
	return &ResponseCompletedMessage{
		ID:   helpers.GenerateUUID(),
		Type: ResponseCompletedMessageType,
	}
}

type ResponseFailedMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
	Code     string      `json:"code"`
	Cause    string      `json:"cause"`
}

func (m *ResponseFailedMessage) GetMessageID() string {
	return m.ID
}

func (m *ResponseFailedMessage) GetMessageType() MessageType {
	return m.Type
}

func (m *ResponseFailedMessage) GetInvokeID() string {
	return m.InvokeID
}

func (m *ResponseFailedMessage) SetInvokeID(invokeID string) {
	m.InvokeID = invokeID
}

func NewResponseFailedMessage(code string, cause string) *ResponseFailedMessage {
	return &ResponseFailedMessage{
		ID:    helpers.GenerateUUID(),
		Type:  ResponseFailedMessageType,
		Cause: cause,
		Code:  code,
	}
}
