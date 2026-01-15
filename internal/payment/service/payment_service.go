package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/transfer"

	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/repository"
)

type PaymentService struct {
	repo       *repository.PaymentRepository
	stripeMode string // "demo" or "real"
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	// Initialize Stripe API key if in real mode
	stripeMode := os.Getenv("STRIPE_MODE")
	if stripeMode == "" {
		stripeMode = "demo" // Default to demo
	}

	if stripeMode == "real" {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	}

	return &PaymentService{
		repo:       repo,
		stripeMode: stripeMode,
	}
}

// CreatePaymentIntent creates a Stripe PaymentIntent and payment record
func (s *PaymentService) CreatePaymentIntent(
	ctx context.Context,
	orderID, userID int64,
	amount float64,
	currency string,
	stripeCustomerID string,
) (*model.Payment, string, error) {
	// Validate amount
	if amount <= 0 {
		return nil, "", fmt.Errorf("amount must be greater than 0")
	}

	// Set default currency
	if currency == "" {
		currency = "USD"
	}

	var stripePaymentIntentID string
	var clientSecret string

	if s.stripeMode == "real" {
		// REAL STRIPE: Create PaymentIntent
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(amount * 100)), // Convert to cents
			Currency: stripe.String(currency),
		}

		if stripeCustomerID != "" {
			params.Customer = stripe.String(stripeCustomerID)
		}

		pi, err := paymentintent.New(params)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create Stripe PaymentIntent: %w", err)
		}

		stripePaymentIntentID = pi.ID
		clientSecret = pi.ClientSecret
	} else {
		// DEMO MODE: Simulate Stripe PaymentIntent
		stripePaymentIntentID = fmt.Sprintf("pi_demo_%d", orderID)
		clientSecret = fmt.Sprintf("pi_demo_%d_secret", orderID)
	}

	// Create payment record
	payment := &model.Payment{
		OrderID:  orderID,
		UserID:   userID,
		Amount:   amount,
		Currency: currency,
		Status:   model.StatusPending,
	}

	// Set Stripe details
	payment.StripePaymentIntentID = sql.NullString{
		String: stripePaymentIntentID,
		Valid:  true,
	}

	if stripeCustomerID != "" {
		payment.StripeCustomerID = sql.NullString{
			String: stripeCustomerID,
			Valid:  true,
		}
	}

	// Create in database
	createdPayment, err := s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create payment: %w", err)
	}

	return createdPayment, clientSecret, nil
}

// ConfirmPayment confirms a payment (called after Stripe confirms)
func (s *PaymentService) ConfirmPayment(
	ctx context.Context,
	paymentID int64,
	stripePaymentIntentID string,
) (*model.Payment, error) {
	// Get payment
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// Verify payment intent ID matches
	if payment.StripePaymentIntentID.String != stripePaymentIntentID {
		return nil, fmt.Errorf("payment intent ID mismatch")
	}

	var chargeID, paymentMethod, cardLast4, cardBrand string

	if s.stripeMode == "real" {
		// REAL STRIPE: Retrieve PaymentIntent to verify status
		pi, err := paymentintent.Get(stripePaymentIntentID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve payment intent: %w", err)
		}

		if pi.Status != stripe.PaymentIntentStatusSucceeded {
			return nil, fmt.Errorf("payment not succeeded: %s", pi.Status)
		}

		// Extract charge details from LatestCharge
		if pi.LatestCharge != nil {
			chargeID = pi.LatestCharge.ID

			if pi.LatestCharge.PaymentMethodDetails != nil && pi.LatestCharge.PaymentMethodDetails.Card != nil {
				paymentMethod = "card"
				cardLast4 = pi.LatestCharge.PaymentMethodDetails.Card.Last4
				cardBrand = string(pi.LatestCharge.PaymentMethodDetails.Card.Brand)
			}
		}
	} else {
		// DEMO MODE: Simulate successful payment
		chargeID = fmt.Sprintf("ch_demo_%d", paymentID)
		paymentMethod = "card"
		cardLast4 = "4242"
		cardBrand = "visa"
	}

	// Update payment with charge details
	err = s.repo.UpdatePaymentWithCharge(
		ctx, paymentID,
		chargeID, paymentMethod, cardLast4, cardBrand,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Update order status to 'paid' via Order Service
	// orderClient.MarkAsPaid(ctx, payment.OrderID)

	return s.repo.GetPaymentByID(ctx, paymentID)
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID int64) (*model.Payment, error) {
	return s.repo.GetPaymentByID(ctx, paymentID)
}

// GetPaymentByOrderID retrieves payment for an order
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	return s.repo.GetPaymentByOrderID(ctx, orderID)
}

// ListPayments lists payments for a user
func (s *PaymentService) ListPayments(
	ctx context.Context,
	userID int64,
	status string,
	page, pageSize int32,
) ([]*model.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.ListPaymentsByUser(ctx, userID, status, page, pageSize)
}

// CreateRefund creates a refund for a payment
func (s *PaymentService) CreateRefund(
	ctx context.Context,
	paymentID int64,
	amount float64,
	reason string,
) (string, error) {
	// Get payment
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return "", err
	}

	// Validate payment is refundable
	if !payment.IsRefundable() {
		return "", fmt.Errorf("payment is not refundable (status: %s)", payment.Status)
	}

	// If amount not specified, refund full remaining amount
	if amount <= 0 {
		amount = payment.CanBeRefundedAmount()
	}

	// Validate refund amount
	if amount > payment.CanBeRefundedAmount() {
		return "", fmt.Errorf("refund amount exceeds available amount")
	}

	var stripeRefundID string

	if s.stripeMode == "real" {
		// REAL STRIPE: Create refund
		if !payment.StripeChargeID.Valid {
			return "", fmt.Errorf("no Stripe charge ID found")
		}

		params := &stripe.RefundParams{
			Charge: stripe.String(payment.StripeChargeID.String),
			Amount: stripe.Int64(int64(amount * 100)),
		}
		if reason != "" {
			params.Reason = stripe.String(reason)
		}

		r, err := refund.New(params)
		if err != nil {
			return "", fmt.Errorf("failed to create Stripe refund: %w", err)
		}

		stripeRefundID = r.ID
	} else {
		// DEMO MODE: Simulate Stripe refund
		stripeRefundID = fmt.Sprintf("re_demo_%d", paymentID)
	}

	// Update payment with refund details
	err = s.repo.UpdatePaymentRefund(ctx, paymentID, amount, reason)
	if err != nil {
		return "", err
	}

	// TODO: If fully refunded, cancel order via Order Service
	// if amount >= payment.CanBeRefundedAmount() {
	//     orderClient.CancelOrder(ctx, payment.OrderID, "Payment refunded")
	// }

	return stripeRefundID, nil
}

// GetRefundStatus retrieves refund information for a payment
func (s *PaymentService) GetRefundStatus(ctx context.Context, paymentID int64) (*model.Payment, error) {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.RefundedAmount == 0 {
		return nil, fmt.Errorf("no refunds for this payment")
	}

	return payment, nil
}

// ========== Payout Methods ==========

// CreatePayout creates a payout for a seller
func (s *PaymentService) CreatePayout(
	ctx context.Context,
	orderID, sellerID, paymentID int64,
	amount float64,
	stripeAccountID string,
) (*model.Payout, error) {
	// Validate amount
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Validate Stripe account ID
	if stripeAccountID == "" {
		return nil, fmt.Errorf("stripe account ID is required")
	}

	var stripeTransferID string

	if s.stripeMode == "real" {
		// REAL STRIPE: Create Transfer (Stripe Connect)
		params := &stripe.TransferParams{
			Amount:      stripe.Int64(int64(amount * 100)),
			Currency:    stripe.String("usd"),
			Destination: stripe.String(stripeAccountID),
		}

		t, err := transfer.New(params)
		if err != nil {
			return nil, fmt.Errorf("failed to create Stripe transfer: %w", err)
		}

		stripeTransferID = t.ID
	} else {
		// DEMO MODE: Simulate transfer
		stripeTransferID = "tr_demo_immediate"
	}

	// Create payout record
	payout := &model.Payout{
		OrderID:   orderID,
		SellerID:  sellerID,
		PaymentID: paymentID,
		Amount:    amount,
		Currency:  "USD",
		Status:    model.PayoutStatusPending,
	}

	payout.StripeAccountID = sql.NullString{
		String: stripeAccountID,
		Valid:  true,
	}

	// Create in database
	createdPayout, err := s.repo.CreatePayout(ctx, payout)
	if err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	// Update status to paid (in real app, this would be updated via webhook)
	err = s.repo.UpdatePayoutStatus(ctx, createdPayout.ID, model.PayoutStatusPaid, stripeTransferID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetPayoutByID(ctx, createdPayout.ID)
}

// GetPayout retrieves a payout by ID
func (s *PaymentService) GetPayout(ctx context.Context, payoutID int64) (*model.Payout, error) {
	return s.repo.GetPayoutByID(ctx, payoutID)
}

// ListPayouts lists payouts for a seller
func (s *PaymentService) ListPayouts(
	ctx context.Context,
	sellerID int64,
	status string,
	page, pageSize int32,
) ([]*model.Payout, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.ListPayoutsBySeller(ctx, sellerID, status, page, pageSize)
}

// HandleStripeWebhook handles Stripe webhook events
func (s *PaymentService) HandleStripeWebhook(ctx context.Context, payload []byte, signature string) (string, error) {
	if s.stripeMode == "demo" {
		// DEMO MODE: Just return success
		return "payment_intent.succeeded", nil
	}

	// REAL STRIPE: Verify webhook signature
	// webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	// event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	// if err != nil {
	//     return "", fmt.Errorf("webhook signature verification failed: %w", err)
	// }

	// TODO: Handle different event types:
	// - payment_intent.succeeded → Update payment status
	// - payment_intent.payment_failed → Mark payment as failed
	// - charge.refunded → Update refund status
	// - transfer.created → Update payout status
	// - transfer.failed → Mark payout as failed

	// For now, return success
	return "payment_intent.succeeded", nil
}

// GetStripeMode returns current Stripe mode (demo or real)
func (s *PaymentService) GetStripeMode() string {
	return s.stripeMode
}
