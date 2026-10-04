package messages

type Message interface {
	GetMessageType() MessageType
	GetMessageID() string
	GetInvokeID() string
	SetInvokeID(invokeID string)
}
type ListMessage []Message

type AnyMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

func (a *AnyMessage) GetMessageID() string {
	return a.ID
}

func (a *AnyMessage) GetMessageType() MessageType {
	return a.Type
}

func (a *AnyMessage) GetInvokeID() string {
	return a.InvokeID
}

func (a *AnyMessage) SetInvokeID(invokeID string) {
	a.InvokeID = invokeID
}
