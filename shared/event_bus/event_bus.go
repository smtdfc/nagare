package event_bus

import (
	"context"
	"sync"
)

type Event[T any] struct {
	Topic   string
	Payload T
}

type EventBus[T any] struct {
	mu          sync.RWMutex
	subscribers map[string][]chan T
}

func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{
		subscribers: make(map[string][]chan T),
	}
}

func (b *EventBus[T]) Subscribe(topic string) (<-chan T, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 100)
	b.subscribers[topic] = append(b.subscribers[topic], ch)

	unsub := func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		subs := b.subscribers[topic]
		for i, sub := range subs {
			if sub == ch {
				b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}

	return ch, unsub
}

func (b *EventBus[T]) Publish(ctx context.Context, topic string, payload T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[topic]; ok {
		for _, ch := range subs {
			select {
			case ch <- payload:
			default:
			}
		}
	}
}
