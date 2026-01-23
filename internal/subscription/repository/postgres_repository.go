package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/model"
)

// PostgresSubscriptionRepository implements SubscriptionRepository with PostgreSQL
type PostgresSubscriptionRepository struct {
	db *pgxpool.Pool
}

// NewPostgresSubscriptionRepository creates a new PostgreSQL repository
func NewPostgresSubscriptionRepository(db *pgxpool.Pool) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{
		db: db,
	}
}

// ============================
// SUBSCRIPTION PLANS
// ============================

func (r *PostgresSubscriptionRepository) GetPlan(ctx context.Context, planID int64) (*model.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly,
		       buyer_fee_percent, seller_fee_percent, features,
		       max_active_listings, max_monthly_transactions,
		       is_active, sort_order, stripe_price_id_monthly, stripe_price_id_yearly,
		       created_at, updated_at
		FROM subscription_plans
		WHERE id = $1
	`

	var plan model.SubscriptionPlan
	err := r.db.QueryRow(ctx, query, planID).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
		&plan.PriceMonthly, &plan.PriceYearly,
		&plan.BuyerFeePercent, &plan.SellerFeePercent, &plan.Features,
		&plan.MaxActiveListings, &plan.MaxMonthlyTransactions,
		&plan.IsActive, &plan.SortOrder, &plan.StripePriceIDMonthly, &plan.StripePriceIDYearly,
		&plan.CreatedAt, &plan.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("plan not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return &plan, nil
}

func (r *PostgresSubscriptionRepository) GetPlanByName(ctx context.Context, name string) (*model.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly,
		       buyer_fee_percent, seller_fee_percent, features,
		       max_active_listings, max_monthly_transactions,
		       is_active, sort_order, stripe_price_id_monthly, stripe_price_id_yearly,
		       created_at, updated_at
		FROM subscription_plans
		WHERE name = $1
	`

	var plan model.SubscriptionPlan
	err := r.db.QueryRow(ctx, query, name).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
		&plan.PriceMonthly, &plan.PriceYearly,
		&plan.BuyerFeePercent, &plan.SellerFeePercent, &plan.Features,
		&plan.MaxActiveListings, &plan.MaxMonthlyTransactions,
		&plan.IsActive, &plan.SortOrder, &plan.StripePriceIDMonthly, &plan.StripePriceIDYearly,
		&plan.CreatedAt, &plan.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("plan '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get plan by name: %w", err)
	}

	return &plan, nil
}

func (r *PostgresSubscriptionRepository) ListPlans(ctx context.Context) ([]*model.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly,
		       buyer_fee_percent, seller_fee_percent, features,
		       max_active_listings, max_monthly_transactions,
		       is_active, sort_order, stripe_price_id_monthly, stripe_price_id_yearly,
		       created_at, updated_at
		FROM subscription_plans
		WHERE is_active = true
		ORDER BY sort_order ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	defer rows.Close()

	var plans []*model.SubscriptionPlan
	for rows.Next() {
		var plan model.SubscriptionPlan
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
			&plan.PriceMonthly, &plan.PriceYearly,
			&plan.BuyerFeePercent, &plan.SellerFeePercent, &plan.Features,
			&plan.MaxActiveListings, &plan.MaxMonthlyTransactions,
			&plan.IsActive, &plan.SortOrder, &plan.StripePriceIDMonthly, &plan.StripePriceIDYearly,
			&plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	return plans, nil
}

func (r *PostgresSubscriptionRepository) ListAllPlans(ctx context.Context) ([]*model.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly,
		       buyer_fee_percent, seller_fee_percent, features,
		       max_active_listings, max_monthly_transactions,
		       is_active, sort_order, stripe_price_id_monthly, stripe_price_id_yearly,
		       created_at, updated_at
		FROM subscription_plans
		ORDER BY sort_order ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all plans: %w", err)
	}
	defer rows.Close()

	var plans []*model.SubscriptionPlan
	for rows.Next() {
		var plan model.SubscriptionPlan
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
			&plan.PriceMonthly, &plan.PriceYearly,
			&plan.BuyerFeePercent, &plan.SellerFeePercent, &plan.Features,
			&plan.MaxActiveListings, &plan.MaxMonthlyTransactions,
			&plan.IsActive, &plan.SortOrder, &plan.StripePriceIDMonthly, &plan.StripePriceIDYearly,
			&plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	return plans, nil
}

func (r *PostgresSubscriptionRepository) CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	query := `
		INSERT INTO subscription_plans (
			name, display_name, description, price_monthly, price_yearly,
			buyer_fee_percent, seller_fee_percent, features,
			max_active_listings, max_monthly_transactions,
			is_active, sort_order, stripe_price_id_monthly, stripe_price_id_yearly
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		plan.Name, plan.DisplayName, plan.Description,
		plan.PriceMonthly, plan.PriceYearly,
		plan.BuyerFeePercent, plan.SellerFeePercent, plan.Features,
		plan.MaxActiveListings, plan.MaxMonthlyTransactions,
		plan.IsActive, plan.SortOrder,
		plan.StripePriceIDMonthly, plan.StripePriceIDYearly,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) UpdatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	query := `
		UPDATE subscription_plans
		SET display_name = $2, description = $3,
		    price_monthly = $4, price_yearly = $5,
		    buyer_fee_percent = $6, seller_fee_percent = $7,
		    features = $8, max_active_listings = $9, max_monthly_transactions = $10,
		    is_active = $11, sort_order = $12,
		    stripe_price_id_monthly = $13, stripe_price_id_yearly = $14
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		plan.ID, plan.DisplayName, plan.Description,
		plan.PriceMonthly, plan.PriceYearly,
		plan.BuyerFeePercent, plan.SellerFeePercent,
		plan.Features, plan.MaxActiveListings, plan.MaxMonthlyTransactions,
		plan.IsActive, plan.SortOrder,
		plan.StripePriceIDMonthly, plan.StripePriceIDYearly,
	)

	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

// ============================
// USER SUBSCRIPTIONS
// ============================

func (r *PostgresSubscriptionRepository) GetUserSubscription(ctx context.Context, subscriptionID int64) (*model.UserSubscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, billing_cycle,
		       current_period_start, current_period_end, cancel_at_period_end,
		       stripe_subscription_id, stripe_customer_id,
		       trial_start, trial_end, metadata,
		       created_at, updated_at, cancelled_at
		FROM user_subscriptions
		WHERE id = $1
	`

	var sub model.UserSubscription
	err := r.db.QueryRow(ctx, query, subscriptionID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
		&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
		&sub.StripeSubscriptionID, &sub.StripeCustomerID,
		&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

func (r *PostgresSubscriptionRepository) GetUserActiveSubscription(ctx context.Context, userID int64) (*model.UserSubscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, billing_cycle,
		       current_period_start, current_period_end, cancel_at_period_end,
		       stripe_subscription_id, stripe_customer_id,
		       trial_start, trial_end, metadata,
		       created_at, updated_at, cancelled_at
		FROM user_subscriptions
		WHERE user_id = $1 AND status = 'active'
		LIMIT 1
	`

	var sub model.UserSubscription
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
		&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
		&sub.StripeSubscriptionID, &sub.StripeCustomerID,
		&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active subscription found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}

	return &sub, nil
}

func (r *PostgresSubscriptionRepository) GetUserActiveSubscriptionWithPlan(ctx context.Context, userID int64) (*model.UserSubscriptionWithPlan, error) {
	query := `
		SELECT us.id, us.user_id, us.plan_id, us.status, us.billing_cycle,
		       us.current_period_start, us.current_period_end, us.cancel_at_period_end,
		       us.stripe_subscription_id, us.stripe_customer_id,
		       us.trial_start, us.trial_end, us.metadata,
		       us.created_at, us.updated_at, us.cancelled_at,
		       sp.id, sp.name, sp.display_name, sp.description,
		       sp.price_monthly, sp.price_yearly,
		       sp.buyer_fee_percent, sp.seller_fee_percent, sp.features,
		       sp.max_active_listings, sp.max_monthly_transactions,
		       sp.is_active, sp.sort_order,
		       sp.stripe_price_id_monthly, sp.stripe_price_id_yearly,
		       sp.created_at, sp.updated_at
		FROM user_subscriptions us
		JOIN subscription_plans sp ON sp.id = us.plan_id
		WHERE us.user_id = $1 AND us.status = 'active'
		LIMIT 1
	`

	var sub model.UserSubscription
	var plan model.SubscriptionPlan

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
		&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
		&sub.StripeSubscriptionID, &sub.StripeCustomerID,
		&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
		&plan.PriceMonthly, &plan.PriceYearly,
		&plan.BuyerFeePercent, &plan.SellerFeePercent, &plan.Features,
		&plan.MaxActiveListings, &plan.MaxMonthlyTransactions,
		&plan.IsActive, &plan.SortOrder,
		&plan.StripePriceIDMonthly, &plan.StripePriceIDYearly,
		&plan.CreatedAt, &plan.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active subscription found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to get subscription with plan: %w", err)
	}

	return &model.UserSubscriptionWithPlan{
		UserSubscription: &sub,
		Plan:             &plan,
	}, nil
}

func (r *PostgresSubscriptionRepository) GetUserSubscriptionHistory(ctx context.Context, userID int64) ([]*model.UserSubscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, billing_cycle,
		       current_period_start, current_period_end, cancel_at_period_end,
		       stripe_subscription_id, stripe_customer_id,
		       trial_start, trial_end, metadata,
		       created_at, updated_at, cancelled_at
		FROM user_subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription history: %w", err)
	}
	defer rows.Close()

	var subscriptions []*model.UserSubscription
	for rows.Next() {
		var sub model.UserSubscription
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
			&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
			&sub.StripeSubscriptionID, &sub.StripeCustomerID,
			&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
			&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *PostgresSubscriptionRepository) CreateUserSubscription(ctx context.Context, subscription *model.UserSubscription) error {
	query := `
		INSERT INTO user_subscriptions (
			user_id, plan_id, status, billing_cycle,
			current_period_start, current_period_end, cancel_at_period_end,
			stripe_subscription_id, stripe_customer_id,
			trial_start, trial_end, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		subscription.UserID, subscription.PlanID, subscription.Status, subscription.BillingCycle,
		subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd, subscription.CancelAtPeriodEnd,
		subscription.StripeSubscriptionID, subscription.StripeCustomerID,
		subscription.TrialStart, subscription.TrialEnd, subscription.Metadata,
	).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) UpdateUserSubscription(ctx context.Context, subscription *model.UserSubscription) error {
	query := `
		UPDATE user_subscriptions
		SET plan_id = $2, status = $3, billing_cycle = $4,
		    current_period_start = $5, current_period_end = $6, cancel_at_period_end = $7,
		    stripe_subscription_id = $8, stripe_customer_id = $9,
		    trial_start = $10, trial_end = $11, metadata = $12
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		subscription.ID, subscription.PlanID, subscription.Status, subscription.BillingCycle,
		subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd, subscription.CancelAtPeriodEnd,
		subscription.StripeSubscriptionID, subscription.StripeCustomerID,
		subscription.TrialStart, subscription.TrialEnd, subscription.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) CancelUserSubscription(ctx context.Context, subscriptionID int64, cancelAtPeriodEnd bool) error {
	query := `
		UPDATE user_subscriptions
		SET status = $2, cancel_at_period_end = $3, cancelled_at = $4
		WHERE id = $1
	`

	status := model.StatusActive
	if !cancelAtPeriodEnd {
		status = model.StatusCancelled
	}

	_, err := r.db.Exec(ctx, query, subscriptionID, status, cancelAtPeriodEnd, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) ReactivateUserSubscription(ctx context.Context, subscriptionID int64) error {
	query := `
		UPDATE user_subscriptions
		SET status = $2, cancel_at_period_end = false, cancelled_at = NULL
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, subscriptionID, model.StatusActive)
	if err != nil {
		return fmt.Errorf("failed to reactivate subscription: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) ExpireUserSubscription(ctx context.Context, subscriptionID int64) error {
	query := `UPDATE user_subscriptions SET status = $2 WHERE id = $1`

	_, err := r.db.Exec(ctx, query, subscriptionID, model.StatusExpired)
	if err != nil {
		return fmt.Errorf("failed to expire subscription: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) GetSubscriptionByStripeID(ctx context.Context, stripeSubscriptionID string) (*model.UserSubscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, billing_cycle,
		       current_period_start, current_period_end, cancel_at_period_end,
		       stripe_subscription_id, stripe_customer_id,
		       trial_start, trial_end, metadata,
		       created_at, updated_at, cancelled_at
		FROM user_subscriptions
		WHERE stripe_subscription_id = $1
	`

	var sub model.UserSubscription
	err := r.db.QueryRow(ctx, query, stripeSubscriptionID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
		&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
		&sub.StripeSubscriptionID, &sub.StripeCustomerID,
		&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription not found for Stripe ID: %s", stripeSubscriptionID)
		}
		return nil, fmt.Errorf("failed to get subscription by Stripe ID: %w", err)
	}

	return &sub, nil
}

func (r *PostgresSubscriptionRepository) CountActiveSubscriptionsByPlan(ctx context.Context, planID int64) (int, error) {
	query := `SELECT COUNT(*) FROM user_subscriptions WHERE plan_id = $1 AND status = 'active'`

	var count int
	err := r.db.QueryRow(ctx, query, planID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	return count, nil
}

func (r *PostgresSubscriptionRepository) GetExpiredSubscriptions(ctx context.Context, beforeDate time.Time) ([]*model.UserSubscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, billing_cycle,
		       current_period_start, current_period_end, cancel_at_period_end,
		       stripe_subscription_id, stripe_customer_id,
		       trial_start, trial_end, metadata,
		       created_at, updated_at, cancelled_at
		FROM user_subscriptions
		WHERE status = 'active' AND current_period_end < $1
	`

	rows, err := r.db.Query(ctx, query, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*model.UserSubscription
	for rows.Next() {
		var sub model.UserSubscription
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
			&sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd,
			&sub.StripeSubscriptionID, &sub.StripeCustomerID,
			&sub.TrialStart, &sub.TrialEnd, &sub.Metadata,
			&sub.CreatedAt, &sub.UpdatedAt, &sub.CancelledAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

// ============================
// SUBSCRIPTION TRANSACTIONS
// ============================

func (r *PostgresSubscriptionRepository) GetTransaction(ctx context.Context, transactionID int64) (*model.SubscriptionTransaction, error) {
	query := `
		SELECT id, user_subscription_id, user_id, plan_id,
		       amount, currency, transaction_type, status,
		       stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id,
		       description, metadata, created_at, processed_at
		FROM subscription_transactions
		WHERE id = $1
	`

	var tx model.SubscriptionTransaction
	err := r.db.QueryRow(ctx, query, transactionID).Scan(
		&tx.ID, &tx.UserSubscriptionID, &tx.UserID, &tx.PlanID,
		&tx.Amount, &tx.Currency, &tx.TransactionType, &tx.Status,
		&tx.StripePaymentIntentID, &tx.StripeInvoiceID, &tx.StripeChargeID,
		&tx.Description, &tx.Metadata, &tx.CreatedAt, &tx.ProcessedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

func (r *PostgresSubscriptionRepository) GetTransactionByStripePaymentIntent(ctx context.Context, paymentIntentID string) (*model.SubscriptionTransaction, error) {
	query := `
		SELECT id, user_subscription_id, user_id, plan_id,
		       amount, currency, transaction_type, status,
		       stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id,
		       description, metadata, created_at, processed_at
		FROM subscription_transactions
		WHERE stripe_payment_intent_id = $1
	`

	var tx model.SubscriptionTransaction
	err := r.db.QueryRow(ctx, query, paymentIntentID).Scan(
		&tx.ID, &tx.UserSubscriptionID, &tx.UserID, &tx.PlanID,
		&tx.Amount, &tx.Currency, &tx.TransactionType, &tx.Status,
		&tx.StripePaymentIntentID, &tx.StripeInvoiceID, &tx.StripeChargeID,
		&tx.Description, &tx.Metadata, &tx.CreatedAt, &tx.ProcessedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found for payment intent: %s", paymentIntentID)
		}
		return nil, fmt.Errorf("failed to get transaction by payment intent: %w", err)
	}

	return &tx, nil
}

func (r *PostgresSubscriptionRepository) GetUserTransactions(ctx context.Context, userID int64, limit, offset int) ([]*model.SubscriptionTransaction, error) {
	query := `
		SELECT id, user_subscription_id, user_id, plan_id,
		       amount, currency, transaction_type, status,
		       stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id,
		       description, metadata, created_at, processed_at
		FROM subscription_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.SubscriptionTransaction
	for rows.Next() {
		var tx model.SubscriptionTransaction
		err := rows.Scan(
			&tx.ID, &tx.UserSubscriptionID, &tx.UserID, &tx.PlanID,
			&tx.Amount, &tx.Currency, &tx.TransactionType, &tx.Status,
			&tx.StripePaymentIntentID, &tx.StripeInvoiceID, &tx.StripeChargeID,
			&tx.Description, &tx.Metadata, &tx.CreatedAt, &tx.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (r *PostgresSubscriptionRepository) GetSubscriptionTransactions(ctx context.Context, subscriptionID int64) ([]*model.SubscriptionTransaction, error) {
	query := `
		SELECT id, user_subscription_id, user_id, plan_id,
		       amount, currency, transaction_type, status,
		       stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id,
		       description, metadata, created_at, processed_at
		FROM subscription_transactions
		WHERE user_subscription_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.SubscriptionTransaction
	for rows.Next() {
		var tx model.SubscriptionTransaction
		err := rows.Scan(
			&tx.ID, &tx.UserSubscriptionID, &tx.UserID, &tx.PlanID,
			&tx.Amount, &tx.Currency, &tx.TransactionType, &tx.Status,
			&tx.StripePaymentIntentID, &tx.StripeInvoiceID, &tx.StripeChargeID,
			&tx.Description, &tx.Metadata, &tx.CreatedAt, &tx.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (r *PostgresSubscriptionRepository) CreateTransaction(ctx context.Context, transaction *model.SubscriptionTransaction) error {
	query := `
		INSERT INTO subscription_transactions (
			user_subscription_id, user_id, plan_id,
			amount, currency, transaction_type, status,
			stripe_payment_intent_id, stripe_invoice_id, stripe_charge_id,
			description, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		transaction.UserSubscriptionID, transaction.UserID, transaction.PlanID,
		transaction.Amount, transaction.Currency, transaction.TransactionType, transaction.Status,
		transaction.StripePaymentIntentID, transaction.StripeInvoiceID, transaction.StripeChargeID,
		transaction.Description, transaction.Metadata,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) UpdateTransaction(ctx context.Context, transaction *model.SubscriptionTransaction) error {
	query := `
		UPDATE subscription_transactions
		SET status = $2, stripe_payment_intent_id = $3, stripe_invoice_id = $4,
		    stripe_charge_id = $5, description = $6, metadata = $7, processed_at = $8
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		transaction.ID, transaction.Status,
		transaction.StripePaymentIntentID, transaction.StripeInvoiceID, transaction.StripeChargeID,
		transaction.Description, transaction.Metadata, transaction.ProcessedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) UpdateTransactionStatus(ctx context.Context, transactionID int64, status model.TransactionStatus) error {
	query := `UPDATE subscription_transactions SET status = $2, processed_at = $3 WHERE id = $1`

	_, err := r.db.Exec(ctx, query, transactionID, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

func (r *PostgresSubscriptionRepository) GetRevenueByPlan(ctx context.Context, startDate, endDate time.Time) (map[int64]float64, error) {
	query := `
		SELECT plan_id, SUM(amount) as revenue
		FROM subscription_transactions
		WHERE status = 'succeeded'
		  AND transaction_type = 'payment'
		  AND created_at BETWEEN $1 AND $2
		GROUP BY plan_id
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by plan: %w", err)
	}
	defer rows.Close()

	revenue := make(map[int64]float64)
	for rows.Next() {
		var planID int64
		var amount float64
		if err := rows.Scan(&planID, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan revenue: %w", err)
		}
		revenue[planID] = amount
	}

	return revenue, nil
}

func (r *PostgresSubscriptionRepository) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) as total_revenue
		FROM subscription_transactions
		WHERE status = 'succeeded'
		  AND transaction_type = 'payment'
		  AND created_at BETWEEN $1 AND $2
	`

	var revenue float64
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&revenue)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}

	return revenue, nil
}
