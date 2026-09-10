package message

type Message interface {
	GetKind() Kind
}
type ListMessage []Message

type ReadOnlyChannel <-chan Message
type WriteOnlyChannel chan<- Message
type Channel chan Message
