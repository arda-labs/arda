package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/realtime"
)

type sseEvent struct {
	Event          string    `json:"event"`
	RecipientType  string    `json:"recipientType"`
	RecipientID    string    `json:"recipientId"`
	NotificationID string    `json:"notificationId,omitempty"`
	UnreadCount    int       `json:"unreadCount,omitempty"`
	OccurredAt     time.Time `json:"occurredAt"`
}

func (s *NotificationService) StreamInAppNotifications(w http.ResponseWriter, r *http.Request) {
	if s.hub == nil {
		http.Error(w, "notification stream unavailable", http.StatusServiceUnavailable)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	recipientType := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("recipient_type")))
	if recipientType == "" {
		recipientType = "USER"
	}
	recipientID := strings.TrimSpace(r.URL.Query().Get("recipient_id"))
	if recipientID == "" {
		http.Error(w, "recipient_id is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	events := s.hub.Subscribe()
	defer s.hub.Unsubscribe(events)

	writeSSE(w, flusher, "heartbeat", sseEvent{Event: "heartbeat", RecipientType: recipientType, RecipientID: recipientID, OccurredAt: time.Now().UTC()})
	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-heartbeat.C:
			writeSSE(w, flusher, "heartbeat", sseEvent{Event: "heartbeat", RecipientType: recipientType, RecipientID: recipientID, OccurredAt: time.Now().UTC()})
		case event := <-events:
			if !matchRecipient(event, recipientType, recipientID) {
				continue
			}
			writeSSE(w, flusher, event.Type, toSSEEvent(event))
		}
	}
}

func matchRecipient(event realtime.Event, recipientType, recipientID string) bool {
	return strings.EqualFold(event.RecipientType, recipientType) && event.RecipientID == recipientID
}

func toSSEEvent(event realtime.Event) sseEvent {
	return sseEvent{
		Event:          event.Type,
		RecipientType:  event.RecipientType,
		RecipientID:    event.RecipientID,
		NotificationID: event.NotificationID,
		UnreadCount:    event.UnreadCount,
		OccurredAt:     event.OccurredAt,
	}
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, eventName string, payload sseEvent) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "event: %s\n", eventName)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}
