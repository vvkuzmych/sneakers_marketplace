package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/subscription/repository"
)

// StripeService handles Stripe-related operations for subscriptions
type StripeService struct {
	subscriptionRepo repository.SubscriptionRepository
	stripeAPIKey     string
	webhookSecret    string
}

// NewStripeService creates a new Stripe service
func NewStripeService(
	repo repository.SubscriptionRepository,
	apiKey string,
	webhookSecret string,
) *StripeService {
	// Set Stripe API key globally
	stripe.Key = apiKey

	return &StripeService{
		subscriptionRepo: repo,
		stripeAPIKey:     apiKey,
		webhookSecret:    webhookSecret,
	}
}

// logWithContext logs a message with context values (request ID, user ID, etc.)
func logWithContext(ctx context.Context, format string, args ...interface{}) {
	// Extract values from context
	requestID := ctx.Value("request_id")
	userID := ctx.Value("user_id")

	// Build prefix with context info
	prefix := "[Stripe]"
	if requestID != nil {
		prefix += fmt.Sprintf(" [req:%v]", requestID)
	}
	if userID != nil {
		prefix += fmt.Sprintf(" [user:%v]", userID)
	}

	// Log with prefix
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s %s", prefix, msg)
}

// logErrorWithContext logs an error with context
func logErrorWithContext(ctx context.Context, err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logWithContext(ctx, "ERROR: %s: %v", msg, err)
}

// ==================== CUSTOMER MANAGEMENT ====================

// CreateStripeCustomer creates a Stripe customer for a user
func (s *StripeService) CreateStripeCustomer(ctx context.Context, userID int64, email string) (string, error) {
	logWithContext(ctx, "Creating Stripe customer for user %d (email: %s)", userID, email)

	params := &stripe.CustomerParams{
		Params: stripe.Params{
			Context: ctx, // Pass context to Stripe API
		},
		Email: stripe.String(email),
		Metadata: map[string]string{
			"user_id": strconv.FormatInt(userID, 10),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		logErrorWithContext(ctx, err, "Failed to create Stripe customer for user %d", userID)
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	logWithContext(ctx, "Successfully created Stripe customer %s for user %d", cust.ID, userID)
	return cust.ID, nil
}

// GetOrCreateStripeCustomer gets existing Stripe customer or creates a new one
func (s *StripeService) GetOrCreateStripeCustomer(ctx context.Context, userID int64, email string) (string, error) {
	// Check if user already has a subscription with Stripe customer ID
	subscription, err := s.subscriptionRepo.GetUserActiveSubscription(ctx, userID)
	if err == nil && subscription.StripeCustomerID != "" {
		return subscription.StripeCustomerID, nil
	}

	// Create new customer
	return s.CreateStripeCustomer(ctx, userID, email)
}

// ==================== SUBSCRIPTION MANAGEMENT ====================

// CreateStripeSubscription creates a Stripe subscription
func (s *StripeService) CreateStripeSubscription(
	ctx context.Context,
	customerID string,
	priceID string,
) (string, error) {
	logWithContext(ctx, "Creating Stripe subscription for customer %s with price %s", customerID, priceID)

	params := &stripe.SubscriptionParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		logErrorWithContext(ctx, err, "Failed to create Stripe subscription for customer %s", customerID)
		return "", fmt.Errorf("failed to create Stripe subscription: %w", err)
	}

	logWithContext(ctx, "Successfully created Stripe subscription %s for customer %s", sub.ID, customerID)
	return sub.ID, nil
}

// CancelStripeSubscription cancels a Stripe subscription
func (s *StripeService) CancelStripeSubscription(ctx context.Context, stripeSubscriptionID string) error {
	params := &stripe.SubscriptionCancelParams{
		Params: stripe.Params{
			Context: ctx,
		},
	}

	_, err := subscription.Cancel(stripeSubscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
	}

	return nil
}

// UpdateStripeSubscription updates a Stripe subscription (e.g., change plan)
func (s *StripeService) UpdateStripeSubscription(
	ctx context.Context,
	stripeSubscriptionID string,
	newPriceID string,
) error {
	// Get current subscription to find the subscription item ID
	sub, err := subscription.Get(stripeSubscriptionID, &stripe.SubscriptionParams{
		Params: stripe.Params{
			Context: ctx,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get Stripe subscription: %w", err)
	}

	if len(sub.Items.Data) == 0 {
		return errors.New("subscription has no items")
	}

	// Update the subscription item with new price
	params := &stripe.SubscriptionParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(sub.Items.Data[0].ID),
				Price: stripe.String(newPriceID),
			},
		},
		ProrationBehavior: stripe.String("create_prorations"),
	}

	_, err = subscription.Update(stripeSubscriptionID, params)
	if err != nil {
		return fmt.Errorf("failed to update Stripe subscription: %w", err)
	}

	return nil
}

// ==================== PAYMENT INTENT ====================

// CreatePaymentIntent creates a payment intent for one-time subscription payment
func (s *StripeService) CreatePaymentIntent(
	ctx context.Context,
	userID int64,
	amount float64,
	currency string,
	metadata map[string]string,
) (string, string, error) {
	// Convert amount to cents
	amountInt := int64(amount * 100)

	// Add user_id to metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["user_id"] = strconv.FormatInt(userID, 10)

	params := &stripe.PaymentIntentParams{
		Params: stripe.Params{
			Context: ctx,
		},
		Amount:   stripe.Int64(amountInt),
		Currency: stripe.String(currency),
		Metadata: metadata,
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return "", "", fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent.ID, intent.ClientSecret, nil
}

// ConfirmPaymentIntent confirms a payment intent
func (s *StripeService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) error {
	params := &stripe.PaymentIntentConfirmParams{
		Params: stripe.Params{
			Context: ctx,
		},
	}

	_, err := paymentintent.Confirm(paymentIntentID, params)
	if err != nil {
		return fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	return nil
}

// ==================== WEBHOOK HANDLING ====================

// HandleWebhook processes Stripe webhook events
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	logWithContext(ctx, "Processing Stripe webhook (payload size: %d bytes)", len(payload))

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		logErrorWithContext(ctx, err, "Webhook signature verification failed")
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	logWithContext(ctx, "Received Stripe event: %s (ID: %s)", event.Type, event.ID)

	// Handle different event types
	switch event.Type {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(ctx, event)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(ctx, event)
	default:
		// Log unhandled event
		logWithContext(ctx, "Unhandled event type: %s", event.Type)
		return nil
	}
}

// handleSubscriptionCreated handles subscription.created webhook
func (s *StripeService) handleSubscriptionCreated(ctx context.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Get our subscription from database by Stripe ID
	dbSub, err := s.subscriptionRepo.GetSubscriptionByStripeID(ctx, subscription.ID)
	if err != nil {
		// Subscription might not exist in our DB yet (race condition)
		fmt.Printf("Subscription not found in DB: %s\n", subscription.ID)
		return nil
	}

	// Update status if needed
	status := mapStripeStatus(subscription.Status)
	if dbSub.Status != status {
		dbSub.Status = status
		if err := s.subscriptionRepo.UpdateUserSubscription(ctx, dbSub); err != nil {
			return fmt.Errorf("failed to update subscription status: %w", err)
		}
	}

	return nil
}

// handleSubscriptionUpdated handles subscription.updated webhook
func (s *StripeService) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Get our subscription from database
	dbSub, err := s.subscriptionRepo.GetSubscriptionByStripeID(ctx, subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Update subscription details
	dbSub.Status = mapStripeStatus(subscription.Status)
	dbSub.CurrentPeriodStart = time.Unix(subscription.CurrentPeriodStart, 0)
	dbSub.CurrentPeriodEnd = time.Unix(subscription.CurrentPeriodEnd, 0)
	dbSub.CancelAtPeriodEnd = subscription.CancelAtPeriodEnd

	if err := s.subscriptionRepo.UpdateUserSubscription(ctx, dbSub); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// handleSubscriptionDeleted handles subscription.deleted webhook
func (s *StripeService) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Get our subscription from database
	dbSub, err := s.subscriptionRepo.GetSubscriptionByStripeID(ctx, subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Mark as canceled
	if err := s.subscriptionRepo.CancelUserSubscription(ctx, dbSub.ID, false); err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// handlePaymentSucceeded handles invoice.payment_succeeded webhook
func (s *StripeService) handlePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	// Get subscription from database
	dbSub, err := s.subscriptionRepo.GetSubscriptionByStripeID(ctx, invoice.Subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Record successful transaction
	transaction := &model.SubscriptionTransaction{
		UserID:                dbSub.UserID,
		UserSubscriptionID:    dbSub.ID,
		PlanID:                dbSub.PlanID,
		TransactionType:       model.TransactionPayment,
		Status:                model.TransactionSucceeded,
		Amount:                float64(invoice.AmountPaid) / 100, // Convert from cents
		Currency:              string(invoice.Currency),
		StripePaymentIntentID: invoice.PaymentIntent.ID,
		StripeInvoiceID:       invoice.ID,
		Description:           "Subscription payment",
		CreatedAt:             time.Now(),
	}

	if err := s.subscriptionRepo.CreateTransaction(ctx, transaction); err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: Failed to record transaction: %v\n", err)
	}

	return nil
}

// handlePaymentFailed handles invoice.payment_failed webhook
func (s *StripeService) handlePaymentFailed(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	// Get subscription from database
	dbSub, err := s.subscriptionRepo.GetSubscriptionByStripeID(ctx, invoice.Subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	// Update subscription status to past_due
	dbSub.Status = model.StatusPastDue
	if err := s.subscriptionRepo.UpdateUserSubscription(ctx, dbSub); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// Record failed transaction
	transaction := &model.SubscriptionTransaction{
		UserID:             dbSub.UserID,
		UserSubscriptionID: dbSub.ID,
		PlanID:             dbSub.PlanID,
		TransactionType:    model.TransactionPayment,
		Status:             model.TransactionFailed,
		Amount:             float64(invoice.AmountDue) / 100,
		Currency:           string(invoice.Currency),
		StripeInvoiceID:    invoice.ID,
		Description:        "Payment failed",
		CreatedAt:          time.Now(),
	}

	if err := s.subscriptionRepo.CreateTransaction(ctx, transaction); err != nil {
		fmt.Printf("Warning: Failed to record failed transaction: %v\n", err)
	}

	// TODO: Send notification to user about payment failure

	return nil
}

// mapStripeStatus maps Stripe subscription status to our status
func mapStripeStatus(stripeStatus stripe.SubscriptionStatus) model.SubscriptionStatus {
	switch stripeStatus {
	case stripe.SubscriptionStatusActive:
		return model.StatusActive
	case stripe.SubscriptionStatusCanceled:
		return model.StatusCancelled
	case stripe.SubscriptionStatusPastDue:
		return model.StatusPastDue
	case stripe.SubscriptionStatusTrialing:
		return model.StatusTrialing
	case stripe.SubscriptionStatusIncomplete, stripe.SubscriptionStatusIncompleteExpired:
		return model.StatusExpired
	default:
		return model.StatusExpired
	}
}

// ==================== COMPLETE FLOW ====================

// SubscribeUserToPlan is the complete flow for subscribing a user to a plan
func (s *StripeService) SubscribeUserToPlan(
	ctx context.Context,
	userID int64,
	email string,
	planID int64,
	billingCycle string,
) (*model.UserSubscription, error) {
	// 1. Get plan details
	plan, err := s.subscriptionRepo.GetPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	// Convert billing cycle
	var cycle model.BillingCycle
	switch billingCycle {
	case "monthly":
		cycle = model.BillingMonthly
	case "yearly":
		cycle = model.BillingYearly
	case "lifetime":
		cycle = model.BillingLifetime
	default:
		return nil, errors.New("invalid billing cycle")
	}

	// Free plan doesn't need Stripe
	if plan.IsFree() {
		return s.createFreeSubscription(ctx, userID, planID, cycle)
	}

	// 2. Get or create Stripe customer
	stripeCustomerID, err := s.GetOrCreateStripeCustomer(ctx, userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// 3. Get Stripe price ID from plan
	var stripePriceID string
	if cycle == model.BillingMonthly {
		stripePriceID = plan.StripePriceIDMonthly
	} else {
		stripePriceID = plan.StripePriceIDYearly
	}

	if stripePriceID == "" {
		return nil, fmt.Errorf("no Stripe price ID configured for plan %s (%s)", plan.Name, billingCycle)
	}

	// 4. Create Stripe subscription
	stripeSubscriptionID, err := s.CreateStripeSubscription(ctx, stripeCustomerID, stripePriceID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe subscription: %w", err)
	}

	// 5. Create subscription in our database
	now := time.Now()
	var endDate time.Time
	if cycle == model.BillingMonthly {
		endDate = now.AddDate(0, 1, 0)
	} else {
		endDate = now.AddDate(1, 0, 0)
	}

	subscription := &model.UserSubscription{
		UserID:               userID,
		PlanID:               planID,
		Status:               model.StatusActive,
		BillingCycle:         cycle,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     endDate,
		CancelAtPeriodEnd:    false,
		StripeCustomerID:     stripeCustomerID,
		StripeSubscriptionID: stripeSubscriptionID,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err = s.subscriptionRepo.CreateUserSubscription(ctx, subscription); err != nil {
		// Rollback: Cancel Stripe subscription
		if cancelErr := s.CancelStripeSubscription(ctx, stripeSubscriptionID); cancelErr != nil {
			fmt.Printf("Warning: Failed to cancel Stripe subscription during rollback: %v\n", cancelErr)
		}
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// 6. Get the created subscription
	created, err := s.subscriptionRepo.GetUserActiveSubscription(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created subscription: %w", err)
	}

	// 7. Record initial transaction (will be confirmed by webhook)
	var amount float64
	if cycle == model.BillingYearly {
		amount = plan.PriceYearly
	} else {
		amount = plan.PriceMonthly
	}

	transaction := &model.SubscriptionTransaction{
		UserID:                userID,
		UserSubscriptionID:    created.ID,
		PlanID:                planID,
		TransactionType:       model.TransactionPayment,
		Status:                model.TransactionPending,
		Amount:                amount,
		Currency:              "USD",
		StripePaymentIntentID: stripeSubscriptionID,
		Description:           fmt.Sprintf("Subscription to %s plan (%s)", plan.DisplayName, billingCycle),
		CreatedAt:             now,
	}

	if err = s.subscriptionRepo.CreateTransaction(ctx, transaction); err != nil {
		fmt.Printf("Warning: Failed to record transaction: %v\n", err)
	}

	return created, nil
}

// createFreeSubscription creates a free subscription without Stripe
func (s *StripeService) createFreeSubscription(
	ctx context.Context,
	userID int64,
	planID int64,
	cycle model.BillingCycle,
) (*model.UserSubscription, error) {
	now := time.Now()
	var endDate time.Time
	if cycle == model.BillingMonthly {
		endDate = now.AddDate(0, 1, 0)
	} else if cycle == model.BillingYearly {
		endDate = now.AddDate(1, 0, 0)
	} else {
		endDate = now.AddDate(100, 0, 0) // Lifetime
	}

	subscription := &model.UserSubscription{
		UserID:             userID,
		PlanID:             planID,
		Status:             model.StatusActive,
		BillingCycle:       cycle,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   endDate,
		CancelAtPeriodEnd:  false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.subscriptionRepo.CreateUserSubscription(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return s.subscriptionRepo.GetUserActiveSubscription(ctx, userID)
}
