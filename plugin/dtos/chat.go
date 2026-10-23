package dtos

const (
	ReceivedChatMessageEvent Event = "ReceivedChatMessageEvent"

	HandshakeEvent        Event = "plugin:handshake"
	HandshakeSuccessEvent Event = "plugin:handshake:success"
	HandshakeFailedEvent  Event = "plugin:handshake:failed"

	PrepareChatSessionEvent        Event = "plugin:prepare_chat_session"
	PrepareChatSessionSuccessEvent Event = "plugin:prepare_chat_session:success"
	PrepareChatSessionFailedEvent  Event = "plugin:prepare_chat_session:failed"

	GetModelProvidersEvent        Event = "plugin:GetModelProvidersEvent"
	GetModelProvidersSuccessEvent Event = "plugin:GetModelProvidersSuccessEvent:success"
	GetModelProvidersFailedEvent  Event = "plugin:GetModelProvidersFailedEvent:failed"

	SendChatMessageEvent        Event = "plugin:send_chat_message"
	SendChatMessageSuccessEvent Event = "plugin:send_chat_message:success"
	SendChatMessageFailedEvent  Event = "plugin:send_chat_message:failed"

	ResetChatChannelEvent        Event = "plugin:reset_chat_channel"
	ResetChatChannelSuccessEvent Event = "plugin:reset_chat_channel:success"
	ResetChatChannelFailedEvent  Event = "plugin:reset_chat_channel:failed"

	SubmitToolApprovalEvent        Event = "plugin:submit_tool_approval"
	SubmitToolApprovalSuccessEvent Event = "plugin:submit_tool_approval:success"
	SubmitToolApprovalFailedEvent  Event = "plugin:submit_tool_approval:failed"
)

type GetModelProvidersEventPayload struct {
}

type GetModelProvidersSuccessEventPayload struct {
	List map[string]string `json:"list"`
}

type GetModelProvidersFailedEventPayload struct {
	Cause string `json:"cause"`
}

type ReceivedChatMessageEventPayload struct {
	SessionID string `json:"sessionID"`
	ChannelID string `json:"channelID"`
	Message   string `json:"message"`
}

type HandshakeEventPayload struct {
	// Deprecated: use RequestID instead
	ID          string `json:"id"`
	PackageName string `json:"packageName"`
	ConnectCode string `json:"connectCode"`
}

type HandshakeSuccessEventPayload struct {
	// Deprecated: use RequestID instead
	ID string `json:"id"`
}

type HandshakeFailedEventPayload struct {
	// Deprecated: use RequestID instead
	ID    string `json:"id"`
	Cause string `json:"cause"`
}

type PrepareChatSessionEventPayload struct {
	ChannelID string `json:"channelID"`
}

type PrepareChatSessionSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
	SessionID string `json:"sessionID"`
}

type PrepareChatSessionFailedEventPayload struct {
	ChannelID string `json:"channelID"`
	Cause     string `json:"cause"`
}

type SendChatMessageEventPayload struct {
	SessionID string `json:"sessionID"`
	Text      string `json:"text"`
}

type SendChatMessageSuccessEventPayload struct {
	SessionID string `json:"sessionID"`
}

type SendChatMessageFailedEventPayload struct {
	Cause string `json:"cause"`
}

type ResetChatChannelEventPayload struct {
	ChannelID string `json:"channelID"`
}

type ResetChatChannelSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
}

type ResetChatChannelFailedEventPayload struct {
	Cause string `json:"cause"`
}

type SubmitToolApprovalEventPayload struct {
	SessionID         string `json:"sessionID"`
	ApprovalRequestID string `json:"approvalRequestID"`
	Accept            bool   `json:"accept"`
}

type SubmitToolApprovalSuccessEventPayload struct{}

type SubmitToolApprovalFailedEventPayload struct {
	Cause string `json:"cause"`
}
