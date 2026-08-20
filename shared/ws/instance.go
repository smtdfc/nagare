package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smtdfc/nagare/shared/dto"
)

type WsInstance struct {
	ID   string
	conn *websocket.Conn
	Auth *dto.AuthPayload
}

func (i *WsInstance) GetConn() *websocket.Conn {
	return i.conn
}

func (i *WsInstance) ReadMessage() (*dto.WsMessage[any], error) {
	_, payload, err := i.conn.ReadMessage()
	if err != nil {
		if errors.Is(err, io.EOF) || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return nil, io.EOF
		}

		return nil, err
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("received empty payload")
	}

	var wsMsg dto.WsMessage[any]
	if err := json.Unmarshal(payload, &wsMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w (Raw: %s)", err, string(payload))
	}

	return &wsMsg, nil
}

func SendMessage[T any](i *WsInstance, event dto.WsEvent, payload T) error {
	raw, err := json.Marshal(&dto.WsMessage[T]{
		Event:   event,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("cannot marshal message: %w", err)
	}

	i.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return i.conn.WriteMessage(websocket.TextMessage, raw)
}

func (i *WsInstance) Close() error {
	if i.conn != nil {
		return i.conn.Close()
	}
	return nil
}
