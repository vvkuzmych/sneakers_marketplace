package model

import "time"

// FeeConfig represents fee configuration for a specific vertical (sneakers, tickets, etc.)
type FeeConfig struct {
	ID                    int64     `json:"id" db:"id"`
	Vertical              string    `json:"vertical" db:"vertical"`
	TransactionFeePercent float64   `json:"transaction_fee_percent" db:"transaction_fee_percent"`
	ProcessingFeeFixed    float64   `json:"processing_fee_fixed" db:"processing_fee_fixed"`
	AuthenticationFee     float64   `json:"authentication_fee" db:"authentication_fee"`
	ShippingBuyerCharge   float64   `json:"shipping_buyer_charge" db:"shipping_buyer_charge"`
	ShippingSellerCost    float64   `json:"shipping_seller_cost" db:"shipping_seller_cost"`
	MinTransactionFee     float64   `json:"min_transaction_fee" db:"min_transaction_fee"`
	MaxTransactionFee     *float64  `json:"max_transaction_fee,omitempty" db:"max_transaction_fee"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// CalculateTransactionFee calculates transaction fee for given price
// Returns fee clamped between min and max (if max is set)
func (fc *FeeConfig) CalculateTransactionFee(price float64) float64 {
	fee := price * (fc.TransactionFeePercent / 100.0)

	// Apply min
	if fee < fc.MinTransactionFee {
		fee = fc.MinTransactionFee
	}

	// Apply max (if set)
	if fc.MaxTransactionFee != nil && fee > *fc.MaxTransactionFee {
		fee = *fc.MaxTransactionFee
	}

	return fee
}

// IsSneakers returns true if this is sneakers vertical
func (fc *FeeConfig) IsSneakers() bool {
	return fc.Vertical == "sneakers"
}

// IsTickets returns true if this is tickets vertical
func (fc *FeeConfig) IsTickets() bool {
	return fc.Vertical == "tickets"
}
