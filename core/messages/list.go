package messages

type Message interface {
	GetType() string
}

type ListMessage []Message
