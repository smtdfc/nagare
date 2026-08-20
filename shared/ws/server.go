package ws

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/smtdfc/nagare/shared/dto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func safeExecute(i *WsInstance, msg *dto.WsMessage[any], cb func(*WsInstance, *dto.WsMessage[any]) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WS Panic Recovered] Event: %s, Error: %v\nStack: %s", msg.Event, r, debug.Stack())
			_ = SendMessage(i, "WS_INTERNAL_ERROR", map[string]string{
				"error": fmt.Sprintf("%v", r),
			})
			err = nil
		}
	}()
	return cb(i, msg)
}

func WsHandler(w http.ResponseWriter, r *http.Request, cb func(*WsInstance, *dto.WsMessage[any]) error) error {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("Upgrade error:", err.Error())
		return err
	}
	defer c.Close()

	i := &WsInstance{ID: uuid.New().String(), conn: c}

	for {
		message, err := i.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				println("WebSocket unexpected closed:", err.Error())
			} else {
				println("WebSocket closed normally.")
			}
			break
		}

		if err := safeExecute(i, message, cb); err != nil {
			break
		}

	}

	return nil
}
