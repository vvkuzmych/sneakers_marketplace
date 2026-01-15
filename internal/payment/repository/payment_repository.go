package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/model"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment creates a new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	query := `
		INSERT INTO payments (
			payment_id, order_id, user_id,
			stripe_payment_intent_id, stripe_customer_id,
			amount, currency, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, query,
		payment.PaymentID, payment.OrderID, payment.UserID,
		payment.StripePaymentIntentID, payment.StripeCustomerID,
		payment.Amount, payment.Currency, payment.Status,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return payment, nil
}

// GetPaymentByID retrieves a payment by ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, paymentID int64) (*model.Payment, error) {
	query := `
		SELECT 
			id, payment_id, order_id, user_id,
			stripe_payment_intent_id, stripe_charge_id, stripe_customer_id,
			amount, currency, status,
			payment_method, card_last4, card_brand,
			refunded_amount, refund_reason,
			processed_at, refunded_at,
			created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment := &model.Payment{}
	err := r.db.QueryRow(ctx, query, paymentID).Scan(
		&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID,
		&payment.StripePaymentIntentID, &payment.StripeChargeID, &payment.StripeCustomerID,
		&payment.Amount, &payment.Currency, &payment.Status,
		&payment.PaymentMethod, &payment.CardLast4, &payment.CardBrand,
		&payment.RefundedAmount, &payment.RefundReason,
		&payment.ProcessedAt, &payment.RefundedAt,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	query := `
		SELECT 
			id, payment_id, order_id, user_id,
			stripe_payment_intent_id, stripe_charge_id, stripe_customer_id,
			amount, currency, status,
			payment_method, card_last4, card_brand,
			refunded_amount, refund_reason,
			processed_at, refunded_at,
			created_at, updated_at
		FROM payments
		WHERE order_id = $1
	`

	payment := &model.Payment{}
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID,
		&payment.StripePaymentIntentID, &payment.StripeChargeID, &payment.StripeCustomerID,
		&payment.Amount, &payment.Currency, &payment.Status,
		&payment.PaymentMethod, &payment.CardLast4, &payment.CardBrand,
		&payment.RefundedAmount, &payment.RefundReason,
		&payment.ProcessedAt, &payment.RefundedAt,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No payment yet
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return payment, nil
}

// UpdatePaymentStatus updates the payment status
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID int64, status string) error {
	query := `
		UPDATE payments
		SET status = $1, processed_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, paymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// UpdatePaymentWithCharge updates payment with Stripe charge details
func (r *PaymentRepository) UpdatePaymentWithCharge(ctx context.Context, paymentID int64, chargeID, paymentMethod, cardLast4, cardBrand string) error {
	query := `
		UPDATE payments
		SET 
			stripe_charge_id = $1,
			payment_method = $2,
			card_last4 = $3,
			card_brand = $4,
			status = $5,
			processed_at = NOW()
		WHERE id = $6
	`

	result, err := r.db.Exec(
		ctx, query,
		chargeID, paymentMethod, cardLast4, cardBrand,
		model.StatusSucceeded, paymentID,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// UpdatePaymentRefund updates payment refund information
func (r *PaymentRepository) UpdatePaymentRefund(ctx context.Context, paymentID int64, refundedAmount float64, reason string) error {
	query := `
		UPDATE payments
		SET 
			refunded_amount = refunded_amount + $1,
			refund_reason = $2,
			refunded_at = NOW(),
			status = CASE 
				WHEN (refunded_amount + $1) >= amount THEN $3
				ELSE $4
			END
		WHERE id = $5
	`

	result, err := r.db.Exec(
		ctx, query,
		refundedAmount, reason,
		model.StatusRefunded, model.StatusPartiallyRefunded,
		paymentID,
	)
	if err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// ListPaymentsByUser lists payments for a user
func (r *PaymentRepository) ListPaymentsByUser(ctx context.Context, userID int64, status string, page, pageSize int32) ([]*model.Payment, int64, error) {
	baseQuery := `
		SELECT 
			id, payment_id, order_id, user_id,
			stripe_payment_intent_id, stripe_charge_id, stripe_customer_id,
			amount, currency, status,
			payment_method, card_last4, card_brand,
			refunded_amount, refund_reason,
			processed_at, refunded_at,
			created_at, updated_at
		FROM payments
		WHERE user_id = $1
	`

	var query string
	var countQuery string
	args := []interface{}{userID}
	argPos := 2

	if status != "" {
		query = baseQuery + " AND status = $" + fmt.Sprintf("%d", argPos)
		countQuery = "SELECT COUNT(*) FROM payments WHERE user_id = $1 AND status = $2"
		args = append(args, status)
		argPos++
	} else {
		query = baseQuery
		countQuery = "SELECT COUNT(*) FROM payments WHERE user_id = $1"
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int64
	countArgs := args[:len(args)-2]
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get payments
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	payments := make([]*model.Payment, 0)
	for rows.Next() {
		payment := &model.Payment{}
		err := rows.Scan(
			&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID,
			&payment.StripePaymentIntentID, &payment.StripeChargeID, &payment.StripeCustomerID,
			&payment.Amount, &payment.Currency, &payment.Status,
			&payment.PaymentMethod, &payment.CardLast4, &payment.CardBrand,
			&payment.RefundedAmount, &payment.RefundReason,
			&payment.ProcessedAt, &payment.RefundedAt,
			&payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	return payments, total, nil
}

// ========== Payout Methods ==========

// CreatePayout creates a new payout
func (r *PaymentRepository) CreatePayout(ctx context.Context, payout *model.Payout) (*model.Payout, error) {
	query := `
		INSERT INTO payouts (
			payout_id, order_id, seller_id, payment_id,
			stripe_account_id, amount, currency, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, query,
		payout.PayoutID, payout.OrderID, payout.SellerID, payout.PaymentID,
		payout.StripeAccountID, payout.Amount, payout.Currency, payout.Status,
	).Scan(&payout.ID, &payout.CreatedAt, &payout.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	return payout, nil
}

// GetPayoutByID retrieves a payout by ID
func (r *PaymentRepository) GetPayoutByID(ctx context.Context, payoutID int64) (*model.Payout, error) {
	query := `
		SELECT 
			id, payout_id, order_id, seller_id, payment_id,
			stripe_transfer_id, stripe_account_id,
			amount, currency, status,
			failure_reason, processed_at,
			created_at, updated_at
		FROM payouts
		WHERE id = $1
	`

	payout := &model.Payout{}
	err := r.db.QueryRow(ctx, query, payoutID).Scan(
		&payout.ID, &payout.PayoutID, &payout.OrderID, &payout.SellerID, &payout.PaymentID,
		&payout.StripeTransferID, &payout.StripeAccountID,
		&payout.Amount, &payout.Currency, &payout.Status,
		&payout.FailureReason, &payout.ProcessedAt,
		&payout.CreatedAt, &payout.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payout not found")
		}
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}

	return payout, nil
}

// UpdatePayoutStatus updates the payout status
func (r *PaymentRepository) UpdatePayoutStatus(ctx context.Context, payoutID int64, status, transferID string) error {
	query := `
		UPDATE payouts
		SET 
			status = $1,
			stripe_transfer_id = $2,
			processed_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, status, transferID, payoutID)
	if err != nil {
		return fmt.Errorf("failed to update payout status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payout not found")
	}

	return nil
}

// ListPayoutsBySeller lists payouts for a seller
func (r *PaymentRepository) ListPayoutsBySeller(ctx context.Context, sellerID int64, status string, page, pageSize int32) ([]*model.Payout, int64, error) {
	baseQuery := `
		SELECT 
			id, payout_id, order_id, seller_id, payment_id,
			stripe_transfer_id, stripe_account_id,
			amount, currency, status,
			failure_reason, processed_at,
			created_at, updated_at
		FROM payouts
		WHERE seller_id = $1
	`

	var query string
	var countQuery string
	args := []interface{}{sellerID}
	argPos := 2

	if status != "" {
		query = baseQuery + " AND status = $" + fmt.Sprintf("%d", argPos)
		countQuery = "SELECT COUNT(*) FROM payouts WHERE seller_id = $1 AND status = $2"
		args = append(args, status)
		argPos++
	} else {
		query = baseQuery
		countQuery = "SELECT COUNT(*) FROM payouts WHERE seller_id = $1"
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int64
	countArgs := args[:len(args)-2]
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payouts: %w", err)
	}

	// Get payouts
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payouts: %w", err)
	}
	defer rows.Close()

	payouts := make([]*model.Payout, 0)
	for rows.Next() {
		payout := &model.Payout{}
		err := rows.Scan(
			&payout.ID, &payout.PayoutID, &payout.OrderID, &payout.SellerID, &payout.PaymentID,
			&payout.StripeTransferID, &payout.StripeAccountID,
			&payout.Amount, &payout.Currency, &payout.Status,
			&payout.FailureReason, &payout.ProcessedAt,
			&payout.CreatedAt, &payout.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payout: %w", err)
		}
		payouts = append(payouts, payout)
	}

	return payouts, total, nil
}
