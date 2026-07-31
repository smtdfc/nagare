package messages

import "github.com/google/uuid"

type ResponseStarted struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func NewResponseStarted() *ResponseStarted {
	return &ResponseStarted{
		ID:   uuid.New().String(),
		Type: RESPONSE_STARTED_MESSAGE,
	}
}

func (r *ResponseStarted) GetType() MessageType {
	return r.Type
}

type ResponseCompleted struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func NewResponseCompleted() *ResponseCompleted {
	return &ResponseCompleted{
		ID:   uuid.New().String(),
		Type: RESPONSE_COMPLETED_MESSAGE,
	}
}

func (r *ResponseCompleted) GetType() MessageType {
	return r.Type
}

type ResponseFailed struct {
	ID    string      `json:"id"`
	Type  MessageType `json:"type"`
	Code  string      `json:"code"`
	Cause string      `json:"cause"`
}

func NewResponseFailed(code string, cause string) *ResponseFailed {
	return &ResponseFailed{
		ID:    uuid.New().String(),
		Type:  RESPONSE_FAILED_MESSAGE,
		Code:  code,
		Cause: cause,
	}
}

func (r *ResponseFailed) GetType() MessageType {
	return r.Type
}
