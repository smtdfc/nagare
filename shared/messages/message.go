package messages

type Message interface {
	GetMessageType() MessageType
	GetMessageID() string
}
type ListMessage []Message

type AnyMessage struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (a *AnyMessage) GetMessageID() string {
	return a.ID
}

func (a *AnyMessage) GetMessageType() MessageType {
	return a.Type
}
