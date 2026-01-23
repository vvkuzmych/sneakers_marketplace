package service

import (
	"context"

	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/repository"
)

// FeeProvider implements the SubscriptionFeeProvider interface from fees package
// It provides fee percentages based on user's active subscription plan
type FeeProvider struct {
	repo repository.SubscriptionRepository
}

// NewFeeProvider creates a new FeeProvider
func NewFeeProvider(repo repository.SubscriptionRepository) *FeeProvider {
	return &FeeProvider{
		repo: repo,
	}
}

// GetUserFeePercentages returns (sellerFee, buyerFee, error) based on user's subscription
// Free:  seller 1%, buyer 1%
// Pro:   seller 0.75%, buyer 1%
// Elite: seller 0.5%, buyer 1%
func (p *FeeProvider) GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error) {
	// Get user's active subscription with plan details
	subscriptionWithPlan, err := p.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		// If no active subscription found, default to Free tier
		return 1.0, 1.0, nil
	}

	// Check subscription status
	if subscriptionWithPlan.Status != "active" && subscriptionWithPlan.Status != "trialing" {
		// Inactive subscription = Free tier
		return 1.0, 1.0, nil
	}

	// Return fee percentages from plan
	// Buyer always pays 1% (fixed), seller pays according to plan tier
	return subscriptionWithPlan.Plan.SellerFeePercent, 1.0, nil
}
