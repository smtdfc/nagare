package rest

import "github.com/smtdfc/nagare/shared/messages"

const (
	SendChatMessageEndpoint   = "/api/v1/user/chat/send"
	CreateChatSessionEndpoint = "/api/v1/user/chat/session/create"
	ListChatSessionsEndpoint  = "/api/v1/user/chat/session/list"
	GetChatSessionEndpoint    = "/api/v1/user/chat/session/:id"
	GetChatHistoryEndpoint    = "/api/v1/user/chat/session/:id/history"
)

type Session struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type SendChatMessageRequest struct {
	SessionID string `json:"sessionID"`
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
