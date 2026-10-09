package tool

type Call struct {
	CallID string
	Name   string
	Args   string
}

type ListToolCall []*Call

func NewToolCall(callID, name, args string) *Call {
	return &Call{
		CallID: callID,
		Name:   name,
		Args:   args,
	}
}
