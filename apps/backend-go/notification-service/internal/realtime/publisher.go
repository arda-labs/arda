package realtime

import "github.com/arda-labs/arda/arda-be-go/services/notification-service/internal/biz"

func (h *Hub) PublishInAppNotification(item *biz.InAppNotification) {
	if item == nil {
		return
	}
	h.Publish(Event{
		Type:           "notification.created",
		RecipientType:  item.RecipientType,
		RecipientID:    item.RecipientID,
		NotificationID: item.ID,
	})
	h.Publish(Event{
		Type:          "notification.count_changed",
		RecipientType: item.RecipientType,
		RecipientID:   item.RecipientID,
	})
}

func (h *Hub) PublishInAppCountChanged(recipientType, recipientID string) {
	h.Publish(Event{
		Type:          "notification.count_changed",
		RecipientType: recipientType,
		RecipientID:   recipientID,
	})
}
