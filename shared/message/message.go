package message

type Message interface {
	GetMessageType() MessageType
}
type ListMessage []Message

type ReadOnlyChannel <-chan Message
type WriteOnlyChannel chan<- Message
type Channel chan Message
