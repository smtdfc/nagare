package messages

type Message interface {
	Kind() string
}

type ListMessage []Message
