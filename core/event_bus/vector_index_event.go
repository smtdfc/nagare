package event_bus

var (
	AddVectorEvent      EventType = "event:vector:add"
	AddVectorErrorEvent EventType = "event:vector:add:error"
)

type AddVectorEventPayload struct {
	RequestID string
	Vector    []float32
	IDs       []uint64
}

func (c *AddVectorEventPayload) GetEventType() EventType {
	return AddVectorEvent
}

type AddVectorErrorEventPayload struct {
	RequestID string
	Error     error
}

func (c *AddVectorErrorEventPayload) GetEventType() EventType {
	return AddVectorErrorEvent
}
