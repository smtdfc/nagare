package contracts

const (
	ReceivedChatMessageEvent Event = "ReceivedChatMessageEvent"

	HandshakeEvent        Event = "plugin:handshake"
	HandshakeSuccessEvent Event = "plugin:handshake:success"
	HandshakeFailedEvent  Event = "plugin:handshake:failed"

	RegisterChatListenerEvent        Event = "plugin:register_chat_listener"
	RegisterChatListenerSuccessEvent Event = "plugin:register_chat_listener:success"
	RegisterChatListenerFailedEvent  Event = "plugin:register_chat_listener:failed"

	PrepareChatChannelEvent        Event = "plugin:prepare_chat_channel"
	PrepareChatChannelSuccessEvent Event = "plugin:prepare_chat_channel:success"
	PrepareChatChannelFailedEvent  Event = "plugin:prepare_chat_channel:failed"

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

type RegisterChatListenerEventPayload struct {
	SessionID string `json:"sessionID"`
}

type RegisterChaListenerSuccessEventPayload struct{}

type RegisterChaListenerFailedEventPayload struct {
	Cause string `json:"cause"`
}

type GetModelProvidersEventPayload struct {
}

type GetModelProvidersSuccessEventPayload struct {
	List map[string]string `json:"list"`
}

type GetModelProvidersFailedEventPayload struct {
	Cause string `json:"cause"`
}

type ReceivedChatMessageEventPayload struct {
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

type PrepareChatChannelEventPayload struct {
	ChannelID string `json:"channelID"`
}

type PrepareChatChannelSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
}

type PrepareChatChannelFailedEventPayload struct {
	ChannelID string `json:"channelID"`
	Cause     string `json:"cause"`
}

type SendChatMessageEventPayload struct {
	ChannelID string `json:"channelID"`
	Text      string `json:"text"`
}

type SendChatMessageSuccessEventPayload struct {
	ChannelID string `json:"channelID"`
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
	ChannelID         string `json:"channelID"`
	ApprovalRequestID string `json:"approvalRequestID"`
	Accept            bool   `json:"accept"`
}

type SubmitToolApprovalSuccessEventPayload struct{}

type SubmitToolApprovalFailedEventPayload struct {
	Cause string `json:"cause"`
}
