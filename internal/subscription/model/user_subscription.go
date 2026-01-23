package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SubscriptionStatus represents subscription status
type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusCancelled SubscriptionStatus = "cancelled"
	StatusExpired   SubscriptionStatus = "expired"
	StatusPastDue   SubscriptionStatus = "past_due"
	StatusTrialing  SubscriptionStatus = "trialing"
)

// BillingCycle represents billing frequency
type BillingCycle string

const (
	BillingMonthly  BillingCycle = "monthly"
	BillingYearly   BillingCycle = "yearly"
	BillingLifetime BillingCycle = "lifetime" // For free plan
)

// UserSubscription represents a user's subscription instance
type UserSubscription struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	PlanID int64 `json:"plan_id"`

	// Status
	Status SubscriptionStatus `json:"status"`

	// Billing
	BillingCycle       BillingCycle `json:"billing_cycle"`
	CurrentPeriodStart time.Time    `json:"current_period_start"`
	CurrentPeriodEnd   time.Time    `json:"current_period_end"`
	CancelAtPeriodEnd  bool         `json:"cancel_at_period_end"`

	// Stripe integration
	StripeSubscriptionID string `json:"stripe_subscription_id,omitempty"`
	StripeCustomerID     string `json:"stripe_customer_id,omitempty"`

	// Trial
	TrialStart *time.Time `json:"trial_start,omitempty"`
	TrialEnd   *time.Time `json:"trial_end,omitempty"`

	// Additional data
	Metadata Metadata `json:"metadata"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

// Metadata is a flexible JSON object for additional data
type Metadata map[string]interface{}

// Value implements driver.Valuer for database storage
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner for database retrieval
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, m)
}

// IsActive checks if subscription is currently active
func (us *UserSubscription) IsActive() bool {
	return us.Status == StatusActive
}

// IsCancelled checks if subscription is cancelled
func (us *UserSubscription) IsCancelled() bool {
	return us.Status == StatusCancelled
}

// IsExpired checks if subscription has expired
func (us *UserSubscription) IsExpired() bool {
	return us.Status == StatusExpired || time.Now().After(us.CurrentPeriodEnd)
}

// IsPastDue checks if payment is past due
func (us *UserSubscription) IsPastDue() bool {
	return us.Status == StatusPastDue
}

// IsOnTrial checks if subscription is in trial period
func (us *UserSubscription) IsOnTrial() bool {
	if us.TrialEnd == nil {
		return false
	}
	return time.Now().Before(*us.TrialEnd)
}

// IsTrialing checks if status is trialing
func (us *UserSubscription) IsTrialing() bool {
	return us.Status == StatusTrialing
}

// DaysUntilExpiration returns days remaining in current period
func (us *UserSubscription) DaysUntilExpiration() int {
	if us.IsExpired() {
		return 0
	}
	duration := time.Until(us.CurrentPeriodEnd)
	return int(duration.Hours() / 24)
}

// IsMonthly checks if billing cycle is monthly
func (us *UserSubscription) IsMonthly() bool {
	return us.BillingCycle == BillingMonthly
}

// IsYearly checks if billing cycle is yearly
func (us *UserSubscription) IsYearly() bool {
	return us.BillingCycle == BillingYearly
}

// IsLifetime checks if billing cycle is lifetime (free plan)
func (us *UserSubscription) IsLifetime() bool {
	return us.BillingCycle == BillingLifetime
}

// WillCancelAtPeriodEnd checks if subscription will auto-cancel
func (us *UserSubscription) WillCancelAtPeriodEnd() bool {
	return us.CancelAtPeriodEnd
}

// GetRenewalDate returns the next billing date
func (us *UserSubscription) GetRenewalDate() time.Time {
	return us.CurrentPeriodEnd
}

// UserSubscriptionWithPlan includes plan details
type UserSubscriptionWithPlan struct {
	*UserSubscription
	Plan *SubscriptionPlan `json:"plan"`
}
