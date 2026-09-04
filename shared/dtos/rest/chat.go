package rest

import "github.com/smtdfc/nagare/shared/message"

const (
	ChatSendMessageEndpoint    = "/api/v1/user/chat/send"
	ChatCreateSessionEndpoint  = "/api/v1/user/chat/session/create"
	ChatGetListSessionEndpoint = "/api/v1/user/chat/session/list"
	ChatGetHistoryEndpoint     = "/api/v1/user/chat/history"
)

type Session struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ChatSendMessageRequest struct {
	SessionID string `json:"session_id"`
	Text      string `json:"text"`
}

type ChatCreateSessionRequest struct {
	Title string `json:"title"`
}

type ChatCreateSessionResponse struct {
	Session *Session `json:"session"`
}

type ChatGetListSessionResponse struct {
	Sessions []*Session `json:"sessions"`
}

type ChatGetHistoryRequest struct {
	SessionID string `json:"session_id"`
}

type ChatGetHistoryResponse struct {
	SessionID string            `json:"session_id"`
	Messages  []message.Message `json:"messages"`
}
