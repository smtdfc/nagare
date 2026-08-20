package bus

import (
	"sync"
)

type Event struct {
	Type    string
	Payload interface{}
}

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
}

func (eb *EventBus) Subscribe(topic string, bufferSize int) (<-chan Event, func()) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, bufferSize)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			eb.mu.Lock()
			defer eb.mu.Unlock()

			conns := eb.subscribers[topic]
			for i, c := range conns {
				if c == ch {
					eb.subscribers[topic] = append(conns[:i], conns[i+1:]...)
					close(ch)
					break
				}
			}
			if len(eb.subscribers[topic]) == 0 {
				delete(eb.subscribers, topic)
			}
		})
	}

	return ch, unsubscribe
}

func (eb *EventBus) Publish(topic string, payload interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if chans, found := eb.subscribers[topic]; found {
		event := Event{
			Type:    topic,
			Payload: payload,
		}

		for _, ch := range chans {
			select {
			case ch <- event:
			default:

			}
		}
	}
}
