package messages

type ResponseCompleted struct {
}

func (r *ResponseCompleted) Kind() string {
	return "ResponseCompleted"
}

type ResponseFailed struct {
	Cause string `json:"cause"`
}

func (r *ResponseFailed) Kind() string {
	return "ResponseFailed"
}
