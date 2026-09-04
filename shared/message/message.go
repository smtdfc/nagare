package message

type Message interface {
	GetKind() MessageKind
}
type ListMessage []Message

type MessageReadOnlyChannel <-chan Message
type MessageWriteOnlyChannel chan<- Message
type MessageChannel chan Message
