package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/olahol/melody"
	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
	"github.com/smtdfc/nagare/shared/helpers"
)

var counter atomic.Int64

type WebsocketCoordinator struct {
	mu          sync.RWMutex
	rooms       map[string]map[*melody.Session]bool
	chatHandler *ChatHandler
}

func (w *WebsocketCoordinator) JoinRoom(roomID string, s *melody.Session) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.rooms[roomID]; !exists {
		w.rooms[roomID] = make(map[*melody.Session]bool)
	}
	w.rooms[roomID][s] = true
}

func (w *WebsocketCoordinator) LeaveRoom(roomID string, s *melody.Session) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if clients, exists := w.rooms[roomID]; exists {
		delete(clients, s)
		if len(clients) == 0 {
			delete(w.rooms, roomID)
		}
	}
}

func (w *WebsocketCoordinator) LeaveAllRooms(s *melody.Session) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for roomID, clients := range w.rooms {
		if _, exists := clients[s]; exists {
			delete(clients, s)
			if len(clients) == 0 {
				delete(w.rooms, roomID)
			}
		}
	}
}

func BroadcastToRoom[T any](w *WebsocketCoordinator, roomID string, event websocket_dtos.WebsocketEvent, data T, exclude *melody.Session) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	clients, exists := w.rooms[roomID]
	if !exists || len(clients) == 0 {
		return nil
	}

	raw, err := helpers.MarshalJson(&websocket_dtos.WebsocketPayload[T]{
		Event: event,
		Data:  data,
	})
	if err != nil {
		return err
	}

	msgBytes := []byte(raw)

	for client := range clients {
		if exclude != nil && client == exclude {
			continue
		}
		_ = client.Write(msgBytes)
	}

	return nil
}
func (w *WebsocketCoordinator) parseMessage(msg []byte) (*websocket_dtos.WebsocketPayload[any], error) {
	return helpers.UnmarshalJson[websocket_dtos.WebsocketPayload[any]](string(msg))
}

func (w *WebsocketCoordinator) HandleMessage(s *melody.Session, msg []byte) {
	message, err := w.parseMessage(msg)
	if err != nil {
		s.Close()
	}

	switch message.Event {
	case websocket_dtos.CHAT_LISTEN_MESSAGE_EVENT:
		w.chatHandler.OnListenMessage(s, w, message)
	}
}

func SendMessage[T any](s *melody.Session, event websocket_dtos.WebsocketEvent, data T) error {
	if s == nil {
		return fmt.Errorf("websocket session is nil")
	}

	raw, err := helpers.MarshalJson(&websocket_dtos.WebsocketPayload[T]{
		Event: event,
		Data:  data,
	})

	if err != nil {
		return fmt.Errorf("failed to marshal websocket payload: %w", err)
	}

	if err := s.Write([]byte(raw)); err != nil {
		return fmt.Errorf("failed to write message to session: %w", err)
	}

	return nil
}

func GetData[T any](payload *websocket_dtos.WebsocketPayload[any]) (*T, error) {
	if payload == nil || payload.Data == nil {
		return nil, fmt.Errorf("payload or data is nil")
	}

	bytesData, err := json.Marshal(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	var result T
	if err := json.Unmarshal(bytesData, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into target type %T: %w", (*T)(nil), err)
	}

	return &result, nil
}

func (w *WebsocketCoordinator) HandleDisconnect(s *melody.Session) {
	w.LeaveAllRooms(s)
}

// @Injectable
func NewWebsocketCoordinator(
	chatHandler *ChatHandler,

) *WebsocketCoordinator {
	return &WebsocketCoordinator{
		rooms:       make(map[string]map[*melody.Session]bool),
		chatHandler: chatHandler,
	}
}
