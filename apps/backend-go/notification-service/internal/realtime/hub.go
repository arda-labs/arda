package realtime

import (
	"sync"
	"time"
)

type Event struct {
	Type           string
	RecipientType  string
	RecipientID    string
	NotificationID string
	UnreadCount    int
	OccurredAt     time.Time
}

type Subscriber chan Event

type Hub struct {
	mu          sync.RWMutex
	subscribers map[Subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{subscribers: map[Subscriber]struct{}{}}
}

func (h *Hub) Subscribe() Subscriber {
	ch := make(Subscriber, 16)
	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch Subscriber) {
	h.mu.Lock()
	if _, ok := h.subscribers[ch]; ok {
		delete(h.subscribers, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *Hub) Publish(event Event) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
