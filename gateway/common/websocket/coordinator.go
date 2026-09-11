package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/olahol/melody"
	websocket_dtos "github.com/smtdfc/nagare/shared/dtos/websocket"
	"github.com/smtdfc/nagare/shared/helpers"
)

type EventHandler func(s *melody.Session, w *Coordinator, payload *websocket_dtos.Payload[any])

type Coordinator struct {
	mu       sync.RWMutex
	rooms    map[string]map[*melody.Session]bool
	handlers map[websocket_dtos.Event]EventHandler
}

func (w *Coordinator) JoinRoom(roomID string, s *melody.Session) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.rooms[roomID]; !exists {
		w.rooms[roomID] = make(map[*melody.Session]bool)
	}
	w.rooms[roomID][s] = true
}

func (w *Coordinator) LeaveRoom(roomID string, s *melody.Session) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if clients, exists := w.rooms[roomID]; exists {
		delete(clients, s)
		if len(clients) == 0 {
			delete(w.rooms, roomID)
		}
	}
}

func (w *Coordinator) LeaveAllRooms(s *melody.Session) {
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

func BroadcastToRoom[T any](w *Coordinator, roomID string, event websocket_dtos.Event, data T, exclude *melody.Session) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	clients, exists := w.rooms[roomID]
	if !exists || len(clients) == 0 {
		return nil
	}

	raw, err := helpers.MarshalJson(&websocket_dtos.Payload[T]{
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

func (w *Coordinator) parseMessage(msg []byte) (*websocket_dtos.Payload[any], error) {
	return helpers.UnmarshalJson[websocket_dtos.Payload[any]](string(msg))
}

func (w *Coordinator) HandleMessage(s *melody.Session, msg []byte) {
	message, err := w.parseMessage(msg)
	if err != nil {
		_ = s.Close()
		return
	}

	w.mu.RLock()
	handler, exists := w.handlers[message.Event]
	w.mu.RUnlock()

	if exists && handler != nil {
		handler(s, w, message)
	}
}

func SendMessage[T any](s *melody.Session, event websocket_dtos.Event, data T) error {
	if s == nil {
		return fmt.Errorf("websocket session is nil")
	}

	raw, err := helpers.MarshalJson(&websocket_dtos.Payload[T]{
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

func GetData[T any](payload *websocket_dtos.Payload[any]) (*T, error) {
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

func (w *Coordinator) HandleDisconnect(s *melody.Session) {
	w.LeaveAllRooms(s)
}

func (w *Coordinator) On(event websocket_dtos.Event, handler EventHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[event] = handler
}

// @Injectable
func NewWebsocketCoordinator() *Coordinator {
	return &Coordinator{
		rooms:    make(map[string]map[*melody.Session]bool),
		handlers: make(map[websocket_dtos.Event]EventHandler),
	}
}
