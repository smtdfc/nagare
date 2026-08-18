package ws

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func Dial(rawURL string) (*WsInstance, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	return &WsInstance{
		ID:   uuid.New().String(),
		conn: conn,
	}, nil
}
