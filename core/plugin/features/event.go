package features

import (
	"sync"
)

type Subscriber struct {
	fn func(interface{})
}

type EventBus struct {
	shards    []*shard
	shardMask uint32
}

type shard struct {
	sync.RWMutex
	subscribers map[string][]func(interface{})
}

func NewEventBus(shardCount uint32) *EventBus {
	eb := &EventBus{
		shards:    make([]*shard, shardCount),
		shardMask: shardCount - 1,
	}
	for i := range eb.shards {
		eb.shards[i] = &shard{subscribers: make(map[string][]func(interface{}))}
	}
	return eb
}

func (eb *EventBus) getShard(topic string) *shard {
	var hash uint32 = 2166136261
	for i := 0; i < len(topic); i++ {
		hash *= 16777619
		hash ^= uint32(topic[i])
	}
	return eb.shards[hash&eb.shardMask]
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	s := eb.getShard(topic)
	s.RLock()
	subs := s.subscribers[topic]
	s.RUnlock()

	for _, sub := range subs {
		go sub(data)
	}
}
