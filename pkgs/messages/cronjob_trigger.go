package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type CronJobTriggerMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

func (c *CronJobTriggerMessage) GetMessageType() MessageType {
	return c.Type
}

func (c *CronJobTriggerMessage) GetMessageID() string {
	return c.ID
}

func (c *CronJobTriggerMessage) GetInvokeID() string {
	return c.InvokeID
}

func (c *CronJobTriggerMessage) SetInvokeID(invokeID string) {
	c.InvokeID = invokeID
}

func NewCronJobTriggerMessage() *CronJobTriggerMessage {
	return &CronJobTriggerMessage{
		ID:   helpers.GenerateUUID(),
		Type: CronJobTriggerMessageType,
	}
}
