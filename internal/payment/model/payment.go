package model

import (
	"database/sql"
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedAt             time.Time      `json:"created_at"`
	RefundedAt            sql.NullTime   `json:"refunded_at"`
	ProcessedAt           sql.NullTime   `json:"processed_at"`
	Currency              string         `json:"currency"`
	PaymentID             string         `json:"payment_id"`
	Status                string         `json:"status"`
	StripePaymentIntentID sql.NullString `json:"stripe_payment_intent_id"`
	StripeCustomerID      sql.NullString `json:"stripe_customer_id"`
	PaymentMethod         sql.NullString `json:"payment_method"`
	CardLast4             sql.NullString `json:"card_last4"`
	CardBrand             sql.NullString `json:"card_brand"`
	RefundReason          sql.NullString `json:"refund_reason"`
	StripeChargeID        sql.NullString `json:"stripe_charge_id"`
	Amount                float64        `json:"amount"`
	RefundedAmount        float64        `json:"refunded_amount"`
	ID                    int64          `json:"id"`
	UserID                int64          `json:"user_id"`
	OrderID               int64          `json:"order_id"`
}

// Payout represents a seller payout
type Payout struct {
	UpdatedAt        time.Time      `json:"updated_at"`
	CreatedAt        time.Time      `json:"created_at"`
	ProcessedAt      sql.NullTime   `json:"processed_at"`
	Currency         string         `json:"currency"`
	PayoutID         string         `json:"payout_id"`
	Status           string         `json:"status"`
	StripeAccountID  sql.NullString `json:"stripe_account_id"`
	StripeTransferID sql.NullString `json:"stripe_transfer_id"`
	FailureReason    sql.NullString `json:"failure_reason"`
	Amount           float64        `json:"amount"`
	ID               int64          `json:"id"`
	PaymentID        int64          `json:"payment_id"`
	SellerID         int64          `json:"seller_id"`
	OrderID          int64          `json:"order_id"`
}

// Payment status constants
const (
	StatusPending           = "pending"
	StatusProcessing        = "processing"
	StatusSucceeded         = "succeeded"
	StatusFailed            = "failed"
	StatusCancelled         = "canceled"
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
