package service

import "context"

// SubscriptionFeeProvider defines interface for getting fee percentages based on user's subscription
type SubscriptionFeeProvider interface {
	// GetUserFeePercentages returns (sellerFee, buyerFee, error)
	// If user has no subscription, returns default free tier fees
	GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error)
}

// DefaultFeeProvider provides default fee percentages when subscription service is not available
type DefaultFeeProvider struct{}

// NewDefaultFeeProvider creates a new default fee provider
func NewDefaultFeeProvider() *DefaultFeeProvider {
	return &DefaultFeeProvider{}
}

// GetUserFeePercentages returns default free tier fees (1% for both)
func (p *DefaultFeeProvider) GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error) {
	// Default free tier fees
	return 1.0, 1.0, nil
}
