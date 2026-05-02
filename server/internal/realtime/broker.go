package realtime

import (
	"encoding/json"
	"sync"

	"opspulse/server/internal/core"
)

type Broker struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

func NewBroker() *Broker {
	return &Broker{clients: make(map[chan []byte]struct{})}
}

func (b *Broker) Subscribe() chan []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan []byte, 8)
	b.clients[ch] = struct{}{}
	return ch
}

func (b *Broker) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, ch)
	close(ch)
}

func (b *Broker) PublishNode(node core.NodeRecord, overview core.Overview) {
	payload, _ := json.Marshal(map[string]any{
		"type":     "node_update",
		"node":     node,
		"overview": overview,
	})
	b.broadcast(payload)
}

func (b *Broker) PublishEvent(event core.Event) {
	payload, _ := json.Marshal(map[string]any{
		"type":  "event",
		"event": event,
	})
	b.broadcast(payload)
}

func (b *Broker) broadcast(msg []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}
