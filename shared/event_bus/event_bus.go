package event_bus

import (
	"context"
	"sync"
)

type Event[T any] struct {
	Topic   string
	Payload T
}

type BaseEventBus[T any] struct {
	mu          sync.RWMutex
	writeMu     sync.Mutex
	subscribers map[string][]chan T
}

func NewBaseEventBus[T any]() *BaseEventBus[T] {
	return &BaseEventBus[T]{
		subscribers: make(map[string][]chan T),
	}
}

func (b *BaseEventBus[T]) Subscribe(eventName string) (<-chan T, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 1000)
	b.subscribers[eventName] = append(b.subscribers[eventName], ch)

	unsub := func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		subs := b.subscribers[eventName]
		for i, sub := range subs {
			if sub == ch {
				b.subscribers[eventName] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}

	return ch, unsub
}
func (b *BaseEventBus[T]) Publish(ctx context.Context, eventName string, payload T) {
	b.mu.RLock()
	subs := b.subscribers[eventName]
	b.mu.RUnlock()

	if len(subs) == 0 {
		return
	}

	b.writeMu.Lock()
	defer b.writeMu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- payload:
		case <-ctx.Done():
			return
		}
	}
}
