package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type TaskTriggerMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

func (t *TaskTriggerMessage) GetMessageType() MessageType {
	return t.Type
}

func (t *TaskTriggerMessage) GetMessageID() string {
	return t.ID
}

func (t *TaskTriggerMessage) GetInvokeID() string {
	return t.InvokeID
}

func (t *TaskTriggerMessage) SetInvokeID(invokeID string) {
	t.InvokeID = invokeID
}

func NewTaskTriggerMessage() *TaskTriggerMessage {
	return &TaskTriggerMessage{
		ID:   helpers.GenerateUUID(),
		Type: TaskTriggerMessageType,
	}
}
