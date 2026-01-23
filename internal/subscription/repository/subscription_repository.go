package repository

import (
	"context"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/model"
)

// SubscriptionRepository defines all database operations for subscriptions
type SubscriptionRepository interface {
	// ============================
	// SUBSCRIPTION PLANS
	// ============================

	// GetPlan retrieves a plan by ID
	GetPlan(ctx context.Context, planID int64) (*model.SubscriptionPlan, error)

	// GetPlanByName retrieves a plan by name ('free', 'pro', 'elite')
	GetPlanByName(ctx context.Context, name string) (*model.SubscriptionPlan, error)

	// ListPlans retrieves all active plans
	ListPlans(ctx context.Context) ([]*model.SubscriptionPlan, error)

	// ListAllPlans retrieves all plans (including inactive)
	ListAllPlans(ctx context.Context) ([]*model.SubscriptionPlan, error)

	// CreatePlan creates a new subscription plan
	CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error

	// UpdatePlan updates an existing plan
	UpdatePlan(ctx context.Context, plan *model.SubscriptionPlan) error

	// ============================
	// USER SUBSCRIPTIONS
	// ============================

	// GetUserSubscription retrieves a subscription by ID
	GetUserSubscription(ctx context.Context, subscriptionID int64) (*model.UserSubscription, error)

	// GetUserActiveSubscription retrieves user's active subscription
	GetUserActiveSubscription(ctx context.Context, userID int64) (*model.UserSubscription, error)

	// GetUserActiveSubscriptionWithPlan retrieves user's active subscription with plan details
	GetUserActiveSubscriptionWithPlan(ctx context.Context, userID int64) (*model.UserSubscriptionWithPlan, error)

	// GetUserSubscriptionHistory retrieves all user subscriptions
	GetUserSubscriptionHistory(ctx context.Context, userID int64) ([]*model.UserSubscription, error)

	// CreateUserSubscription creates a new subscription for a user
	CreateUserSubscription(ctx context.Context, subscription *model.UserSubscription) error

	// UpdateUserSubscription updates an existing subscription
	UpdateUserSubscription(ctx context.Context, subscription *model.UserSubscription) error

	// CancelUserSubscription marks subscription as cancelled
	CancelUserSubscription(ctx context.Context, subscriptionID int64, cancelAtPeriodEnd bool) error

	// ReactivateUserSubscription reactivates a cancelled subscription
	ReactivateUserSubscription(ctx context.Context, subscriptionID int64) error

	// ExpireUserSubscription marks subscription as expired
	ExpireUserSubscription(ctx context.Context, subscriptionID int64) error

	// GetSubscriptionByStripeID retrieves subscription by Stripe subscription ID
	GetSubscriptionByStripeID(ctx context.Context, stripeSubscriptionID string) (*model.UserSubscription, error)

	// CountActiveSubscriptionsByPlan counts active subscriptions for a plan
	CountActiveSubscriptionsByPlan(ctx context.Context, planID int64) (int, error)

	// GetExpiredSubscriptions retrieves subscriptions that need to be expired
	GetExpiredSubscriptions(ctx context.Context, beforeDate time.Time) ([]*model.UserSubscription, error)

	// ============================
	// SUBSCRIPTION TRANSACTIONS
	// ============================

	// GetTransaction retrieves a transaction by ID
	GetTransaction(ctx context.Context, transactionID int64) (*model.SubscriptionTransaction, error)

	// GetTransactionByStripePaymentIntent retrieves transaction by Stripe payment intent ID
	GetTransactionByStripePaymentIntent(ctx context.Context, paymentIntentID string) (*model.SubscriptionTransaction, error)

	// GetUserTransactions retrieves all transactions for a user
	GetUserTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.SubscriptionTransaction, error)

	// GetSubscriptionTransactions retrieves all transactions for a subscription
	GetSubscriptionTransactions(ctx context.Context, subscriptionID int64) ([]*model.SubscriptionTransaction, error)

	// CreateTransaction creates a new transaction
	CreateTransaction(ctx context.Context, transaction *model.SubscriptionTransaction) error

	// UpdateTransaction updates an existing transaction
	UpdateTransaction(ctx context.Context, transaction *model.SubscriptionTransaction) error

	// UpdateTransactionStatus updates transaction status and processed_at
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status model.TransactionStatus) error

	// GetRevenueByPlan calculates total revenue by plan
	GetRevenueByPlan(ctx context.Context, startDate, endDate time.Time) (map[int64]float64, error)

	// GetTotalRevenue calculates total revenue in a date range
	GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}
