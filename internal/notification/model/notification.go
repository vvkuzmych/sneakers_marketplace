package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Notification types
const (
	TypeMatchCreated     = "match_created"
	TypeOrderCreated     = "order_created"
	TypeOrderPaid        = "order_paid"
	TypeOrderShipped     = "order_shipped"
	TypeOrderDelivered   = "order_delivered"
	TypeOrderCompleted   = "order_completed"
	TypePaymentSucceeded = "payment_succeeded"
	TypePaymentFailed    = "payment_failed"
	TypeRefundIssued     = "refund_issued"
	TypePayoutCompleted  = "payout_completed"
)

// Notification represents a notification in the system
type Notification struct {
	CreatedAt   time.Time
	EmailSentAt sql.NullTime
	ReadAt      sql.NullTime
	PushSentAt  sql.NullTime
	Title       string
	Message     string
	Type        string
	Data        json.RawMessage
	ID          int64
	UserID      int64
	EmailSent   bool
	PushSent    bool
	IsRead      bool
}

// NotificationPreferences represents user's notification preferences
type NotificationPreferences struct {
	UpdatedAt            time.Time
	CreatedAt            time.Time
	UserID               int64
	EmailPayoutCompleted bool
	EmailOrderShipped    bool
	EmailPaymentReceived bool
	EmailOrderCreated    bool
	PushEnabled          bool
	PushMatchCreated     bool
	PushOrderUpdates     bool
	PushPaymentUpdates   bool
	InAppEnabled         bool
	EmailMatchCreated    bool
	EmailEnabled         bool
}

// ShouldSendEmail checks if email notification should be sent for this type
func (p *NotificationPreferences) ShouldSendEmail(notificationType string) bool {
	if !p.EmailEnabled {
		return false
	}

	switch notificationType {
	case TypeMatchCreated:
		return p.EmailMatchCreated
	case TypeOrderCreated, TypeOrderPaid, TypeOrderShipped, TypeOrderDelivered, TypeOrderCompleted:
		return p.EmailOrderCreated || p.EmailOrderShipped
	case TypePaymentSucceeded, TypePaymentFailed, TypeRefundIssued:
		return p.EmailPaymentReceived
	case TypePayoutCompleted:
		return p.EmailPayoutCompleted
	default:
		return true // Send by default for unknown types
	}
}

// ShouldSendPush checks if push notification should be sent for this type
func (p *NotificationPreferences) ShouldSendPush(notificationType string) bool {
	if !p.PushEnabled {
		return false
	}

	switch notificationType {
	case TypeMatchCreated:
		return p.PushMatchCreated
	case TypeOrderCreated, TypeOrderPaid, TypeOrderShipped, TypeOrderDelivered, TypeOrderCompleted:
		return p.PushOrderUpdates
	case TypePaymentSucceeded, TypePaymentFailed, TypeRefundIssued, TypePayoutCompleted:
		return p.PushPaymentUpdates
	default:
		return true
	}
}

// NewNotification creates a new notification
func NewNotification(userID int64, notifType, title, message string, data map[string]interface{}) (*Notification, error) {
	var jsonData json.RawMessage
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		jsonData = b
	}

	return &Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      jsonData,
		EmailSent: false,
		PushSent:  false,
		IsRead:    false,
		CreatedAt: time.Now(),
	}, nil
}

// GetDataField extracts a field from the JSON data
func (n *Notification) GetDataField(key string) (interface{}, bool) {
	if len(n.Data) == 0 {
		return nil, false
	}

	var data map[string]interface{}
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return nil, false
	}

	value, ok := data[key]
	return value, ok
}
