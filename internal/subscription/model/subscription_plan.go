package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SubscriptionPlan represents a subscription tier (Free, Pro, Elite)
type SubscriptionPlan struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`         // 'free', 'pro', 'elite'
	DisplayName string `json:"display_name"` // 'Free', 'Pro', 'Elite'
	Description string `json:"description"`

	// Pricing
	PriceMonthly float64 `json:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly"`

	// Fee structure
	BuyerFeePercent  float64 `json:"buyer_fee_percent"`
	SellerFeePercent float64 `json:"seller_fee_percent"`

	// Features and limits
	Features               Features `json:"features"`
	MaxActiveListings      *int     `json:"max_active_listings"`      // NULL = unlimited
	MaxMonthlyTransactions *int     `json:"max_monthly_transactions"` // NULL = unlimited

	// Metadata
	IsActive             bool   `json:"is_active"`
	SortOrder            int    `json:"sort_order"`
	StripePriceIDMonthly string `json:"stripe_price_id_monthly,omitempty"`
	StripePriceIDYearly  string `json:"stripe_price_id_yearly,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Features is a JSON array of plan features
type Features []string

// Value implements driver.Valuer for database storage
func (f Features) Value() (driver.Value, error) {
	if f == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(f)
}

// Scan implements sql.Scanner for database retrieval
func (f *Features) Scan(value interface{}) error {
	if value == nil {
		*f = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, f)
}

// IsFree checks if this is the free plan
func (sp *SubscriptionPlan) IsFree() bool {
	return sp.Name == "free"
}

// IsPro checks if this is the pro plan
func (sp *SubscriptionPlan) IsPro() bool {
	return sp.Name == "pro"
}

// IsElite checks if this is the elite plan
func (sp *SubscriptionPlan) IsElite() bool {
	return sp.Name == "elite"
}

// GetMonthlyPrice returns monthly price as cents (for Stripe)
func (sp *SubscriptionPlan) GetMonthlyPriceCents() int64 {
	return int64(sp.PriceMonthly * 100)
}

// GetYearlyPrice returns yearly price as cents (for Stripe)
func (sp *SubscriptionPlan) GetYearlyPriceCents() int64 {
	return int64(sp.PriceYearly * 100)
}

// HasUnlimitedListings checks if plan has unlimited listings
func (sp *SubscriptionPlan) HasUnlimitedListings() bool {
	return sp.MaxActiveListings == nil
}

// HasUnlimitedTransactions checks if plan has unlimited transactions
func (sp *SubscriptionPlan) HasUnlimitedTransactions() bool {
	return sp.MaxMonthlyTransactions == nil
}

// GetYearlySavings calculates savings when paying yearly vs monthly
func (sp *SubscriptionPlan) GetYearlySavings() float64 {
	monthlyTotal := sp.PriceMonthly * 12
	return monthlyTotal - sp.PriceYearly
}

// GetYearlySavingsPercent calculates savings percentage
func (sp *SubscriptionPlan) GetYearlySavingsPercent() float64 {
	if sp.PriceMonthly == 0 {
		return 0
	}
	savings := sp.GetYearlySavings()
	monthlyTotal := sp.PriceMonthly * 12
	return (savings / monthlyTotal) * 100
}
