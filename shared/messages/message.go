package messages

import (
	"fmt"
)

type Message interface {
	GetType() MessageType
}

type ListMessage []Message

var EMPTY_LIST = ListMessage{}

func PrintMessage(raw Message) {
	switch message := raw.(type) {
	case *Text:
		fmt.Printf("[TEXT] %s:%s\n", message.Role, message.Content)
	case *Reasoning:
		fmt.Printf("[Reasoning] Agent:%s\n", message.Content)
	case *ResponseStarted:
		fmt.Printf("[System] Response started\n")
	case *ResponseCompleted:
		fmt.Printf("[System] Response completed\n")
	case *ResponseFailed:
		fmt.Printf("[System] Response failed. Cause: %s\n", message.Cause)
	}
}
