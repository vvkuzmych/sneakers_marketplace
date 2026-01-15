package model

import (
	"database/sql"
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID        int64  `json:"id"`
	PaymentID string `json:"payment_id"`

	// Relations
	OrderID int64 `json:"order_id"`
	UserID  int64 `json:"user_id"`

	// Stripe details
	StripePaymentIntentID sql.NullString `json:"stripe_payment_intent_id"`
	StripeChargeID        sql.NullString `json:"stripe_charge_id"`
	StripeCustomerID      sql.NullString `json:"stripe_customer_id"`

	// Amount
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	// Status
	Status string `json:"status"`

	// Payment method
	PaymentMethod sql.NullString `json:"payment_method"`
	CardLast4     sql.NullString `json:"card_last4"`
	CardBrand     sql.NullString `json:"card_brand"`

	// Refund
	RefundedAmount float64        `json:"refunded_amount"`
	RefundReason   sql.NullString `json:"refund_reason"`

	// Timestamps
	ProcessedAt sql.NullTime `json:"processed_at"`
	RefundedAt  sql.NullTime `json:"refunded_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Payout represents a seller payout
type Payout struct {
	ID       int64  `json:"id"`
	PayoutID string `json:"payout_id"`

	// Relations
	OrderID   int64 `json:"order_id"`
	SellerID  int64 `json:"seller_id"`
	PaymentID int64 `json:"payment_id"`

	// Stripe Connect
	StripeTransferID sql.NullString `json:"stripe_transfer_id"`
	StripeAccountID  sql.NullString `json:"stripe_account_id"`

	// Amount
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	// Status
	Status string `json:"status"`

	// Failure info
	FailureReason sql.NullString `json:"failure_reason"`

	// Timestamps
	ProcessedAt sql.NullTime `json:"processed_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Payment status constants
const (
	StatusPending           = "pending"
	StatusProcessing        = "processing"
	StatusSucceeded         = "succeeded"
	StatusFailed            = "failed"
	StatusCancelled         = "cancelled"
	StatusRefunded          = "refunded"
	StatusPartiallyRefunded = "partially_refunded"
)

// Payout status constants
const (
	PayoutStatusPending    = "pending"
	PayoutStatusProcessing = "processing"
	PayoutStatusPaid       = "paid"
	PayoutStatusFailed     = "failed"
	PayoutStatusReversed   = "reversed"
)

// IsSuccessful checks if payment was successful
func (p *Payment) IsSuccessful() bool {
	return p.Status == StatusSucceeded ||
		p.Status == StatusPartiallyRefunded
}

// IsRefundable checks if payment can be refunded
func (p *Payment) IsRefundable() bool {
	return p.Status == StatusSucceeded &&
		p.RefundedAmount < p.Amount
}

// CanBeRefundedAmount checks how much can be refunded
func (p *Payment) CanBeRefundedAmount() float64 {
	if !p.IsRefundable() {
		return 0
	}
	return p.Amount - p.RefundedAmount
}
