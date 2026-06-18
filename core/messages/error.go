package messages

type StreamErrorMessage struct {
	Cause string
}

func (m *StreamErrorMessage) GetType() string {
	return "StreamError"
}
