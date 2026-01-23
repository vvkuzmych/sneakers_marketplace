package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/repository"
)

// SubscriptionService handles business logic for subscriptions
type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

// ==================== SUBSCRIPTION PLANS ====================

// GetAllPlans returns all available subscription plans
func (s *SubscriptionService) GetAllPlans(ctx context.Context) ([]*model.SubscriptionPlan, error) {
	return s.repo.ListAllPlans(ctx)
}

// GetPlanByID returns a plan by ID
func (s *SubscriptionService) GetPlanByID(ctx context.Context, planID int64) (*model.SubscriptionPlan, error) {
	return s.repo.GetPlan(ctx, planID)
}

// GetPlanByName returns a plan by name (free, pro, elite)
func (s *SubscriptionService) GetPlanByName(ctx context.Context, name string) (*model.SubscriptionPlan, error) {
	return s.repo.GetPlanByName(ctx, name)
}

// GetActivePlans returns all active subscription plans
func (s *SubscriptionService) GetActivePlans(ctx context.Context) ([]*model.SubscriptionPlan, error) {
	return s.repo.ListPlans(ctx)
}

// ==================== USER SUBSCRIPTIONS ====================

// GetUserSubscription returns the active subscription for a user
func (s *SubscriptionService) GetUserSubscription(ctx context.Context, userID int64) (*model.UserSubscriptionWithPlan, error) {
	return s.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
}

// GetUserSubscriptionHistory returns all subscriptions for a user
func (s *SubscriptionService) GetUserSubscriptionHistory(ctx context.Context, userID int64) ([]*model.UserSubscription, error) {
	return s.repo.GetUserSubscriptionHistory(ctx, userID)
}

// CreateSubscription creates a new subscription for a user
func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID, planID int64, billingCycle string) (*model.UserSubscription, error) {
	// Validate plan exists
	plan, err := s.repo.GetPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	// Convert and validate billing cycle
	var cycle model.BillingCycle
	switch billingCycle {
	case "monthly":
		cycle = model.BillingMonthly
	case "yearly":
		cycle = model.BillingYearly
	case "lifetime":
		cycle = model.BillingLifetime
	default:
		return nil, errors.New("invalid billing cycle: must be 'monthly', 'yearly', or 'lifetime'")
	}

	// Check if user already has an active subscription
	existing, err := s.repo.GetUserActiveSubscription(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing subscription: %w", err)
	}

	if existing != nil && existing.Status == model.StatusActive {
		return nil, errors.New("user already has an active subscription")
	}

	// Calculate dates
	now := time.Now()
	var endDate time.Time
	if cycle == model.BillingMonthly {
		endDate = now.AddDate(0, 1, 0) // Add 1 month
	} else if cycle == model.BillingYearly {
		endDate = now.AddDate(1, 0, 0) // Add 1 year
	} else {
		// Lifetime - set far future
		endDate = now.AddDate(100, 0, 0)
	}

	// Create subscription
	subscription := &model.UserSubscription{
		UserID:             userID,
		PlanID:             planID,
		Status:             model.StatusActive,
		BillingCycle:       cycle,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   endDate,
		CancelAtPeriodEnd:  false, // Don't cancel at period end (auto-renew)
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = s.repo.CreateUserSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	// Return the created subscription
	return s.repo.GetUserActiveSubscription(ctx, userID)
}

// UpgradeSubscription upgrades a user to a higher tier plan
func (s *SubscriptionService) UpgradeSubscription(ctx context.Context, userID, newPlanID int64) (*model.UserSubscription, error) {
	// Get current subscription
	current, err := s.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current subscription: %w", err)
	}

	if current.Status != model.StatusActive {
		return nil, errors.New("current subscription is not active")
	}

	// Get new plan
	newPlan, err := s.repo.GetPlan(ctx, newPlanID)
	if err != nil {
		return nil, fmt.Errorf("new plan not found: %w", err)
	}

	if !newPlan.IsActive {
		return nil, errors.New("new plan is not active")
	}

	// Validate upgrade (can't downgrade using this method)
	if newPlan.PriceMonthly <= current.Plan.PriceMonthly {
		return nil, errors.New("can only upgrade to a higher tier plan")
	}

	// Cancel current subscription
	if err := s.CancelSubscription(ctx, current.ID); err != nil {
		return nil, fmt.Errorf("failed to cancel current subscription: %w", err)
	}

	// Create new subscription
	return s.CreateSubscription(ctx, userID, newPlanID, string(current.BillingCycle))
}

// DowngradeSubscription downgrades a user to a lower tier plan
func (s *SubscriptionService) DowngradeSubscription(ctx context.Context, userID, newPlanID int64) (*model.UserSubscription, error) {
	// Get current subscription
	current, err := s.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current subscription: %w", err)
	}

	if current.Status != model.StatusActive {
		return nil, errors.New("current subscription is not active")
	}

	// Get new plan
	newPlan, err := s.repo.GetPlan(ctx, newPlanID)
	if err != nil {
		return nil, fmt.Errorf("new plan not found: %w", err)
	}

	if !newPlan.IsActive {
		return nil, errors.New("new plan is not active")
	}

	// Validate downgrade
	if newPlan.PriceMonthly >= current.Plan.PriceMonthly {
		return nil, errors.New("can only downgrade to a lower tier plan")
	}

	// Schedule downgrade at end of current billing period
	// Mark subscription to cancel at period end
	if err := s.repo.CancelUserSubscription(ctx, current.ID, true); err != nil {
		return nil, fmt.Errorf("failed to schedule downgrade: %w", err)
	}

	// Return updated subscription
	return s.repo.GetUserSubscription(ctx, current.ID)
}

// CancelSubscription cancels a user's subscription
func (s *SubscriptionService) CancelSubscription(ctx context.Context, subscriptionID int64) error {
	// Get subscription
	subscription, err := s.repo.GetUserSubscription(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	if subscription.Status != model.StatusActive {
		return errors.New("subscription is not active")
	}

	// Cancel subscription immediately
	return s.repo.CancelUserSubscription(ctx, subscriptionID, false)
}

// RenewSubscription renews an expired subscription
func (s *SubscriptionService) RenewSubscription(ctx context.Context, subscriptionID int64) (*model.UserSubscription, error) {
	// Get subscription
	subscription, err := s.repo.GetUserSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	if subscription.Status != model.StatusExpired && subscription.Status != model.StatusCancelled {
		return nil, errors.New("subscription is still active")
	}

	// Calculate new dates
	now := time.Now()
	var endDate time.Time
	if subscription.BillingCycle == model.BillingMonthly {
		endDate = now.AddDate(0, 1, 0)
	} else if subscription.BillingCycle == model.BillingYearly {
		endDate = now.AddDate(1, 0, 0)
	} else {
		endDate = now.AddDate(100, 0, 0) // Lifetime
	}

	// Reactivate subscription
	if err := s.repo.ReactivateUserSubscription(ctx, subscriptionID); err != nil {
		return nil, fmt.Errorf("failed to reactivate subscription: %w", err)
	}

	// Update dates
	subscription.Status = model.StatusActive
	subscription.CurrentPeriodStart = now
	subscription.CurrentPeriodEnd = endDate
	subscription.CancelAtPeriodEnd = false
	subscription.UpdatedAt = now

	if err := s.repo.UpdateUserSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	return s.repo.GetUserSubscription(ctx, subscriptionID)
}

// GetExpiringSubscriptions returns subscriptions expiring within X days
func (s *SubscriptionService) GetExpiringSubscriptions(ctx context.Context, days int) ([]*model.UserSubscription, error) {
	beforeDate := time.Now().AddDate(0, 0, days)
	return s.repo.GetExpiredSubscriptions(ctx, beforeDate)
}

// ProcessExpiredSubscriptions marks expired subscriptions as expired
func (s *SubscriptionService) ProcessExpiredSubscriptions(ctx context.Context) error {
	// Get all subscriptions that should be expired
	expiredSubs, err := s.repo.GetExpiredSubscriptions(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to get expired subscriptions: %w", err)
	}

	// Mark each as expired
	for _, sub := range expiredSubs {
		if err := s.repo.ExpireUserSubscription(ctx, sub.ID); err != nil {
			// Log error but continue with others
			fmt.Printf("Warning: Failed to expire subscription %d: %v\n", sub.ID, err)
		}
	}

	return nil
}

// ==================== TRANSACTIONS ====================

// CreateTransaction records a subscription transaction
func (s *SubscriptionService) CreateTransaction(ctx context.Context, transaction *model.SubscriptionTransaction) error {
	// Validate transaction
	if transaction.UserID == 0 {
		return errors.New("user_id is required")
	}
	if transaction.UserSubscriptionID == 0 {
		return errors.New("user_subscription_id is required")
	}
	if transaction.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if transaction.TransactionType != model.TransactionPayment && transaction.TransactionType != model.TransactionRefund {
		return errors.New("invalid transaction type")
	}

	// Set defaults
	if transaction.Currency == "" {
		transaction.Currency = "USD"
	}
	if transaction.Status == "" {
		transaction.Status = model.TransactionPending
	}

	now := time.Now()
	transaction.CreatedAt = now

	return s.repo.CreateTransaction(ctx, transaction)
}

// GetUserTransactions returns all transactions for a user
func (s *SubscriptionService) GetUserTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.SubscriptionTransaction, error) {
	return s.repo.GetUserTransactions(ctx, userID, limit, offset)
}

// GetSubscriptionTransactions returns all transactions for a subscription
func (s *SubscriptionService) GetSubscriptionTransactions(ctx context.Context, subscriptionID int64) ([]*model.SubscriptionTransaction, error) {
	return s.repo.GetSubscriptionTransactions(ctx, subscriptionID)
}

// UpdateTransactionStatus updates the status of a transaction
func (s *SubscriptionService) UpdateTransactionStatus(ctx context.Context, transactionID int64, status model.TransactionStatus) error {
	// Validate status
	validStatuses := map[model.TransactionStatus]bool{
		model.TransactionPending:   true,
		model.TransactionSucceeded: true,
		model.TransactionFailed:    true,
		model.TransactionRefunded:  true,
	}

	if !validStatuses[status] {
		return errors.New("invalid transaction status")
	}

	return s.repo.UpdateTransactionStatus(ctx, transactionID, status)
}

// ==================== FEE CALCULATION ====================

// GetUserFeePercentages returns the fee percentages for a user based on their subscription
func (s *SubscriptionService) GetUserFeePercentages(ctx context.Context, userID int64) (float64, float64, error) {
	// Get user's subscription
	subscription, err := s.repo.GetUserActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		// If no subscription, return free tier fees
		freePlan, err := s.repo.GetPlanByName(ctx, "free")
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get free plan: %w", err)
		}
		return freePlan.SellerFeePercent, freePlan.BuyerFeePercent, nil
	}

	// Return fees from user's plan
	return subscription.Plan.SellerFeePercent, subscription.Plan.BuyerFeePercent, nil
}

// ==================== STATISTICS ====================

// GetSubscriptionStats returns statistics about subscriptions
func (s *SubscriptionService) GetSubscriptionStats(ctx context.Context) (map[string]interface{}, error) {
	// Get all plans and count active subscriptions for each
	plans, err := s.repo.ListAllPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	countByPlan := make(map[int64]int)
	for _, plan := range plans {
		count, err := s.repo.CountActiveSubscriptionsByPlan(ctx, plan.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to count subscriptions for plan %d: %w", plan.ID, err)
		}
		countByPlan[plan.ID] = count
	}

	// Total revenue
	startDate := time.Time{}
	endDate := time.Now()
	revenueByPlan, err := s.repo.GetRevenueByPlan(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue: %w", err)
	}

	totalRevenue := 0.0
	for _, revenue := range revenueByPlan {
		totalRevenue += revenue
	}

	stats := map[string]interface{}{
		"by_plan":         countByPlan,
		"revenue_by_plan": revenueByPlan,
		"total_revenue":   totalRevenue,
		"generated_at":    time.Now(),
	}

	return stats, nil
}

// GetMonthlyRevenue returns the revenue for a specific month
func (s *SubscriptionService) GetMonthlyRevenue(ctx context.Context, year, month int) (float64, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	revenueByPlan, err := s.repo.GetRevenueByPlan(ctx, start, end)
	if err != nil {
		return 0, err
	}

	totalRevenue := 0.0
	for _, revenue := range revenueByPlan {
		totalRevenue += revenue
	}

	return totalRevenue, nil
}
