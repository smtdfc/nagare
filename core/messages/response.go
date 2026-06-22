package messages

type ResponseCreatedMessage struct{}

func (m *ResponseCreatedMessage) GetType() string {
	return "ResponseCreated"
}

type ResponseCompletedMessage struct {
}

func (m *ResponseCompletedMessage) GetType() string {
	return "ResponseCompleted"
}

type ResponseStatsMessage struct {
	InputTokens     int64
	OutputTokens    int64
	ReasoningTokens int64
	TotalTokens     int64
}

func (m *ResponseStatsMessage) GetType() string {
	return "ResponseStats"
}

type ResponseFailedMessage struct {
	Code    string
	Message string
}

func (m *ResponseFailedMessage) GetType() string {
	return "ResponseFailed"
}
