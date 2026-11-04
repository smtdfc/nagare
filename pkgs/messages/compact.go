package messages

import "github.com/smtdfc/nagare/pkgs/helpers"

type ConversationSummaryMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
	Summary  string      `json:"summary"`
}

// GetInvokeID implements [Message].
func (c *ConversationSummaryMessage) GetInvokeID() string {
	return c.InvokeID
}

// GetMessageID implements [Message].
func (c *ConversationSummaryMessage) GetMessageID() string {
	return c.ID
}

// GetMessageType implements [Message].
func (c *ConversationSummaryMessage) GetMessageType() MessageType {
	return c.Type
}

// SetInvokeID implements [Message].
func (c *ConversationSummaryMessage) SetInvokeID(invokeID string) {
	c.InvokeID = invokeID
}

func NewConversationSummaryMessage(summary string) *ConversationSummaryMessage {
	return &ConversationSummaryMessage{
		ID:      helpers.GenerateUUID(),
		Summary: summary,
		Type:    ConversationSummaryMessageType,
	}
}

type CompactSuccessMessage struct {
	ID       string      `json:"id"`
	Type     MessageType `json:"type"`
	InvokeID string      `json:"invokeID"`
}

// GetInvokeID implements [Message].
func (c *CompactSuccessMessage) GetInvokeID() string {
	return c.InvokeID
}

// GetMessageID implements [Message].
func (c *CompactSuccessMessage) GetMessageID() string {
	return c.ID
}

// GetMessageType implements [Message].
func (c *CompactSuccessMessage) GetMessageType() MessageType {
	return c.Type
}

// SetInvokeID implements [Message].
func (c *CompactSuccessMessage) SetInvokeID(invokeID string) {
	c.InvokeID = invokeID
}

func NewCompactSuccessMessage() *CompactSuccessMessage {
	return &CompactSuccessMessage{
		ID:   helpers.GenerateUUID(),
		Type: CompactSuccessMessageType,
	}
}
