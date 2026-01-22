package model

import (
	"database/sql"
	"time"
)

// Order represents an order created from a matched bid/ask
type Order struct {
	ID          int64  `json:"id"`
	OrderNumber string `json:"order_number"`
	MatchID     int64  `json:"match_id"`

	// Parties
	BuyerID  int64 `json:"buyer_id"`
	SellerID int64 `json:"seller_id"`

	// Product details
	ProductID int64 `json:"product_id"`
	SizeID    int64 `json:"size_id"`

	// Pricing
	Price        float64 `json:"price"`
	Quantity     int32   `json:"quantity"`
	BuyerFee     float64 `json:"buyer_fee"`
	SellerFee    float64 `json:"seller_fee"`
	PlatformFee  float64 `json:"platform_fee"`
	TotalAmount  float64 `json:"total_amount"`  // price + buyer_fee
	SellerPayout float64 `json:"seller_payout"` // price - seller_fee

	// Status
	Status string `json:"status"`

	// Shipping
	ShippingAddressID sql.NullInt64  `json:"shipping_address_id"`
	TrackingNumber    sql.NullString `json:"tracking_number"`
	Carrier           sql.NullString `json:"carrier"`

	// Timestamps
	PaymentAt   sql.NullTime `json:"payment_at"`
	ShippedAt   sql.NullTime `json:"shipped_at"`
	DeliveredAt sql.NullTime `json:"delivered_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
	CancelledAt sql.NullTime `json:"canceled_at"`

	// Notes
	BuyerNotes         sql.NullString `json:"buyer_notes"`
	SellerNotes        sql.NullString `json:"seller_notes"`
	AdminNotes         sql.NullString `json:"admin_notes"`
	CancellationReason sql.NullString `json:"cancellation_reason"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderStatusHistory represents a status change event for an order
type OrderStatusHistory struct {
	ID         int64          `json:"id"`
	OrderID    int64          `json:"order_id"`
	FromStatus sql.NullString `json:"from_status"`
	ToStatus   string         `json:"to_status"`
	Note       sql.NullString `json:"note"`
	CreatedBy  string         `json:"created_by"` // system, buyer, seller, admin
	CreatedAt  time.Time      `json:"created_at"`
}

// Order status constants
const (
	StatusPendingPayment = "pending_payment"
	StatusPaid           = "paid"
	StatusProcessing     = "processing"
	StatusShipped        = "shipped"
	StatusDelivered      = "delivered"
	StatusCompleted      = "completed"
	StatusCancelled      = "canceled"
	StatusRefunded       = "refunded"
)

// ValidStatuses returns all valid order statuses
func ValidStatuses() []string {
	return []string{
		StatusPendingPayment,
		StatusPaid,
		StatusProcessing,
		StatusShipped,
		StatusDelivered,
		StatusCompleted,
		StatusCancelled,
		StatusRefunded,
	}
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

// CanTransitionTo checks if status transition is allowed
func (o *Order) CanTransitionTo(newStatus string) bool {
	// Define allowed transitions
	allowedTransitions := map[string][]string{
		StatusPendingPayment: {StatusPaid, StatusCancelled},
		StatusPaid:           {StatusProcessing, StatusRefunded, StatusCancelled},
		StatusProcessing:     {StatusShipped, StatusCancelled},
		StatusShipped:        {StatusDelivered, StatusCancelled},
		StatusDelivered:      {StatusCompleted, StatusCancelled},
		StatusCompleted:      {StatusRefunded}, // Can refund completed orders
		StatusCancelled:      {},               // Final state
		StatusRefunded:       {},               // Final state
	}

	allowed, exists := allowedTransitions[o.Status]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

// IsFinalStatus checks if the order is in a final state
func (o *Order) IsFinalStatus() bool {
	return o.Status == StatusCompleted ||
		o.Status == StatusCancelled ||
		o.Status == StatusRefunded
}

// CanBeCancelled checks if the order can be canceled
func (o *Order) CanBeCancelled() bool {
	return o.Status != StatusCompleted &&
		o.Status != StatusCancelled &&
		o.Status != StatusRefunded
}
