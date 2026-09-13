package message

type Message interface {
	GetMessageType() MessageType
}
type ListMessage []Message

type ReadOnlyChannel <-chan Message
type WriteOnlyChannel chan<- Message
type Channel chan Message

type AnyMessage struct {
	Type MessageType `json:"type"`
}

func (a *AnyMessage) GetMessageType() MessageType {
	return a.Type
}
