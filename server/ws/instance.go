package ws

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/smtdfc/nagare/shared/dto"
)

type WsInstance struct {
	conn *websocket.Conn
}

func SendMessage[T any](i *WsInstance, event dto.WsEvent, payload T) error {
	raw, err := json.Marshal(&dto.WsMessage[T]{
		Event:   event,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("Cannot send message")
	}

	if err := i.conn.WriteMessage(websocket.BinaryMessage, raw); err != nil {
		return err
	}

	return nil
}

func ReadMessage[T any](i *WsInstance) (*dto.WsMessage[T], error) {
	var (
		mt  int
		raw []byte
		err error
	)

	if mt, raw, err = i.conn.ReadMessage(); err != nil {
		return nil, err
	}

	var msg dto.WsMessage[T]
	err = json.Unmarshal(raw, &msg)
	if err != nil {
		return nil, err
	}

	_ = mt
	return &msg, nil

	return nil, fmt.Errorf("Unsupported message type")
}

func NewWsInstance(c *websocket.Conn) *WsInstance {
	return &WsInstance{
		conn: c,
	}
}
