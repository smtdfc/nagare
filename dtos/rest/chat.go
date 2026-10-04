package rest

import "github.com/smtdfc/nagare/pkgs/messages"

const (
	SendChatMessageEndpoint              = "/api/v1/user/chat/send"
	CreateChatSessionEndpoint            = "/api/v1/user/chat/session/create"
	ListChatSessionsEndpoint             = "/api/v1/user/chat/session/list"
	GetChatSessionEndpoint               = "/api/v1/user/chat/session/:id"
	UpdateChatSessionLLMSettingsEndpoint = "/api/v1/user/chat/session/:id/settings"
	GetChatHistoryEndpoint               = "/api/v1/user/chat/session/:id/history"
	DeleteChatSessionEndpoint            = "/api/v1/user/chat/session/:id"
	ArchiveChatSessionEndpoint           = "/api/v1/user/chat/session/:id/archive"
	DuplicateChatSessionEndpoint         = "/api/v1/user/chat/session/:id/duplicate"
)

type Session struct {
	ID                 string           `json:"id"`
	Title              string           `json:"title"`
	IsArchive          bool             `json:"isArchive"`
	CurrentLLMModel    string           `json:"currentLLMModel,omitempty"`
	CurrentLLMProvider *LLMProviderInfo `json:"currentLLMProvider,omitempty"`
}

type SendChatMessageRequest struct {
	SessionID string `json:"sessionID"`
	Text      string `json:"text"`
}

type SendChatMessageResponse struct {
	InvokeID string `json:"invokeID"`
}

type CreateChatSessionRequest struct {
	Title              string `json:"title"`
	CurrentLLMProvider string `json:"currentLLMProvider"`
	CurrentLLMModel    string `json:"currentLLMModel"`
}

type CreateChatSessionResponse struct {
	Session *Session `json:"session"`
}

type ListChatSessionsResponse struct {
	Sessions []*Session `json:"sessions"`
}

type ListChatSessionsRequest struct {
	Limit  int `json:"limit" query:"limit"`
	Offset int `json:"offset" query:"offset"`
}

type GetChatHistoryRequest struct {
	SessionID string `json:"sessionID"`
	Limit     int    `json:"limit" query:"limit"`
	BeforeID  string `json:"beforeID" query:"beforeID"`
}

type GetChatHistoryResponse struct {
	SessionID  string             `json:"sessionID"`
	Messages   []messages.Message `json:"messages"`
	NextCursor string             `json:"nextCursor"`
}

type GetChatSessionResponse struct {
	Session *Session `json:"session"`
}

type UpdateChatSessionLLMSettingsRequest struct {
	CurrentLLMProvider string `json:"currentLLMProvider"`
	CurrentLLMModel    string `json:"currentLLMModel"`
}

type UpdateChatSessionLLMSettingsResponse struct {
	Session *Session `json:"session"`
}

type ArchiveChatSessionRequest struct {
	IsArchive bool `json:"isArchive"`
}

type ArchiveChatSessionResponse struct {
	Session *Session `json:"session"`
}

type DuplicateChatSessionResponse struct {
	Session *Session `json:"session"`
}
