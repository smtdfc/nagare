package messages

type Message interface {
	GetMessageType() MessageType
}
type ListMessage []Message

type AnyMessage struct {
	Type MessageType `json:"type"`
}

func (a *AnyMessage) GetMessageType() MessageType {
	return a.Type
}
