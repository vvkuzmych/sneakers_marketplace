package model

import (
	"fmt"
	"time"
)

// TransactionType represents the type of subscription transaction
type TransactionType string

const (
	TransactionPayment   TransactionType = "payment"
	TransactionRefund    TransactionType = "refund"
	TransactionUpgrade   TransactionType = "upgrade"
	TransactionDowngrade TransactionType = "downgrade"
)

// TransactionStatus represents transaction status
type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionSucceeded TransactionStatus = "succeeded"
	TransactionFailed    TransactionStatus = "failed"
	TransactionRefunded  TransactionStatus = "refunded"
)

// SubscriptionTransaction represents a payment or billing event
type SubscriptionTransaction struct {
	ID                 int64 `json:"id"`
	UserSubscriptionID int64 `json:"user_subscription_id"`
	UserID             int64 `json:"user_id"`
	PlanID             int64 `json:"plan_id"`

	// Transaction details
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`

	// Type and status
	TransactionType TransactionType   `json:"transaction_type"`
	Status          TransactionStatus `json:"status"`

	// Stripe references
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`
	StripeInvoiceID       string `json:"stripe_invoice_id,omitempty"`
	StripeChargeID        string `json:"stripe_charge_id,omitempty"`

	// Additional info
	Description string   `json:"description,omitempty"`
	Metadata    Metadata `json:"metadata"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

// IsPayment checks if transaction is a payment
func (st *SubscriptionTransaction) IsPayment() bool {
	return st.TransactionType == TransactionPayment
}

// IsRefund checks if transaction is a refund
func (st *SubscriptionTransaction) IsRefund() bool {
	return st.TransactionType == TransactionRefund
}

// IsUpgrade checks if transaction is an upgrade
func (st *SubscriptionTransaction) IsUpgrade() bool {
	return st.TransactionType == TransactionUpgrade
}

// IsDowngrade checks if transaction is a downgrade
func (st *SubscriptionTransaction) IsDowngrade() bool {
	return st.TransactionType == TransactionDowngrade
}

// IsPending checks if transaction is pending
func (st *SubscriptionTransaction) IsPending() bool {
	return st.Status == TransactionPending
}

// IsSucceeded checks if transaction succeeded
func (st *SubscriptionTransaction) IsSucceeded() bool {
	return st.Status == TransactionSucceeded
}

// IsFailed checks if transaction failed
func (st *SubscriptionTransaction) IsFailed() bool {
	return st.Status == TransactionFailed
}

// IsRefunded checks if transaction was refunded
func (st *SubscriptionTransaction) IsRefunded() bool {
	return st.Status == TransactionRefunded
}

// GetAmountCents returns amount in cents (for Stripe)
func (st *SubscriptionTransaction) GetAmountCents() int64 {
	return int64(st.Amount * 100)
}

// GetFormattedAmount returns formatted amount string
func (st *SubscriptionTransaction) GetFormattedAmount() string {
	return formatCurrency(st.Amount, st.Currency)
}

// formatCurrency formats a float64 amount as currency
func formatCurrency(amount float64, currency string) string {
	// Simple formatting for USD
	if currency == "USD" {
		return "$" + formatFloat(amount, 2)
	}
	return formatFloat(amount, 2) + " " + currency
}

// formatFloat formats a float64 to string with fixed decimal places
func formatFloat(val float64, precision int) string {
	switch precision {
	case 2:
		return fmt.Sprintf("%.2f", val)
	default:
		return fmt.Sprintf("%f", val)
	}
}
