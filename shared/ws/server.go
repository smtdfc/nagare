package ws

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/smtdfc/nagare/shared/dto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

		if err := cb(i, message); err != nil {
			println("Callback error:", err.Error())
			break
		}
	}

	return nil
}
