package messages

type ResponseStarted struct {
}

func (r *ResponseStarted) Kind() string {
	return "ResponseStarted"
}

type ResponseCompleted struct {
}

func (r *ResponseCompleted) Kind() string {
	return "ResponseCompleted"
}

type ResponseFailed struct {
	Code  string
	Cause string `json:"cause"`
}

func (r *ResponseFailed) Kind() string {
	return "ResponseFailed"
}
