package rest

import "github.com/smtdfc/nagare/shared/message"

const (
	SendChatMessageEndpoint   = "/api/v1/user/chat/send"
	CreateChatSessionEndpoint = "/api/v1/user/chat/session/create"
	ListChatSessionsEndpoint  = "/api/v1/user/chat/session/list"
	GetChatHistoryEndpoint    = "/api/v1/user/chat/history"
)

type Session struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type SendChatMessageRequest struct {
	SessionID string `json:"session_id"`
	Text      string `json:"text"`
}

type CreateChatSessionRequest struct {
	Title string `json:"title"`
}

type CreateChatSessionResponse struct {
	Session *Session `json:"session"`
}

type ListChatSessionsResponse struct {
	Sessions []*Session `json:"sessions"`
}

type GetChatHistoryRequest struct {
	SessionID string `json:"session_id"`
}

type GetChatHistoryResponse struct {
	SessionID string            `json:"session_id"`
	Messages  []message.Message `json:"messages"`
}
