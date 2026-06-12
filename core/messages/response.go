package messages

type ResponseCreatedMessage struct{}

func (m *ResponseCreatedMessage) GetType() string {
	return "ResponseCreated"
}

type ResponseCompletedMessage struct{}

func (m *ResponseCompletedMessage) GetType() string {
	return "ResponseCompleted"
}

type ResponseFailedMessage struct {
	Code    string
	Message string
}

func (m *ResponseFailedMessage) GetType() string {
	return "ResponseFailed"
}
