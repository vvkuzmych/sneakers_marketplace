package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/order/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order
func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	query := `
		INSERT INTO orders (
			order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			buyer_notes
		) VALUES (
			$1, $2,
			$3, $4,
			$5, $6,
			$7, $8,
			$9, $10, $11,
			$12, $13,
			$14,
			$15,
			$16
		)
		RETURNING id, order_number, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, query,
		order.OrderNumber, order.MatchID,
		order.BuyerID, order.SellerID,
		order.ProductID, order.SizeID,
		order.Price, order.Quantity,
		order.BuyerFee, order.SellerFee, order.PlatformFee,
		order.TotalAmount, order.SellerPayout,
		order.Status,
		order.ShippingAddressID,
		order.BuyerNotes,
	).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	query := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.ID, &order.OrderNumber, &order.MatchID,
		&order.BuyerID, &order.SellerID,
		&order.ProductID, &order.SizeID,
		&order.Price, &order.Quantity,
		&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
		&order.TotalAmount, &order.SellerPayout,
		&order.Status,
		&order.ShippingAddressID,
		&order.TrackingNumber, &order.Carrier,
		&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
		&order.CompletedAt, &order.CancelledAt,
		&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
		&order.CancellationReason,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// GetOrderByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	query := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
		WHERE order_number = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(ctx, query, orderNumber).Scan(
		&order.ID, &order.OrderNumber, &order.MatchID,
		&order.BuyerID, &order.SellerID,
		&order.ProductID, &order.SizeID,
		&order.Price, &order.Quantity,
		&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
		&order.TotalAmount, &order.SellerPayout,
		&order.Status,
		&order.ShippingAddressID,
		&order.TrackingNumber, &order.Carrier,
		&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
		&order.CompletedAt, &order.CancelledAt,
		&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
		&order.CancellationReason,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// ListOrders retrieves orders with optional status filter and pagination
func (r *OrderRepository) ListOrders(ctx context.Context, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	baseQuery := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
	`

	var query string
	var countQuery string
	var args []interface{}
	argPos := 1

	if status != "" {
		query = baseQuery + " WHERE status = $" + fmt.Sprintf("%d", argPos)
		countQuery = "SELECT COUNT(*) FROM orders WHERE status = $1"
		args = append(args, status)
		argPos++
	} else {
		query = baseQuery
		countQuery = "SELECT COUNT(*) FROM orders"
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int64
	countArgs := args[:len(args)-2] // Exclude LIMIT and OFFSET
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get orders
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*model.Order, 0)
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.MatchID,
			&order.BuyerID, &order.SellerID,
			&order.ProductID, &order.SizeID,
			&order.Price, &order.Quantity,
			&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
			&order.TotalAmount, &order.SellerPayout,
			&order.Status,
			&order.ShippingAddressID,
			&order.TrackingNumber, &order.Carrier,
			&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
			&order.CompletedAt, &order.CancelledAt,
			&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
			&order.CancellationReason,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// GetBuyerOrders retrieves orders for a specific buyer
func (r *OrderRepository) GetBuyerOrders(ctx context.Context, buyerID int64, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	baseQuery := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
		WHERE buyer_id = $1
	`

	var query string
	var countQuery string
	args := []interface{}{buyerID}
	argPos := 2

	if status != "" {
		query = baseQuery + " AND status = $" + fmt.Sprintf("%d", argPos)
		countQuery = "SELECT COUNT(*) FROM orders WHERE buyer_id = $1 AND status = $2"
		args = append(args, status)
		argPos++
	} else {
		query = baseQuery
		countQuery = "SELECT COUNT(*) FROM orders WHERE buyer_id = $1"
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int64
	countArgs := args[:len(args)-2] // Exclude LIMIT and OFFSET
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count buyer orders: %w", err)
	}

	// Get orders
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list buyer orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*model.Order, 0)
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.MatchID,
			&order.BuyerID, &order.SellerID,
			&order.ProductID, &order.SizeID,
			&order.Price, &order.Quantity,
			&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
			&order.TotalAmount, &order.SellerPayout,
			&order.Status,
			&order.ShippingAddressID,
			&order.TrackingNumber, &order.Carrier,
			&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
			&order.CompletedAt, &order.CancelledAt,
			&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
			&order.CancellationReason,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// GetSellerOrders retrieves orders for a specific seller
func (r *OrderRepository) GetSellerOrders(ctx context.Context, sellerID int64, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	baseQuery := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
		WHERE seller_id = $1
	`

	var query string
	var countQuery string
	args := []interface{}{sellerID}
	argPos := 2

	if status != "" {
		query = baseQuery + " AND status = $" + fmt.Sprintf("%d", argPos)
		countQuery = "SELECT COUNT(*) FROM orders WHERE seller_id = $1 AND status = $2"
		args = append(args, status)
		argPos++
	} else {
		query = baseQuery
		countQuery = "SELECT COUNT(*) FROM orders WHERE seller_id = $1"
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int64
	countArgs := args[:len(args)-2] // Exclude LIMIT and OFFSET
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seller orders: %w", err)
	}

	// Get orders
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list seller orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*model.Order, 0)
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.MatchID,
			&order.BuyerID, &order.SellerID,
			&order.ProductID, &order.SizeID,
			&order.Price, &order.Quantity,
			&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
			&order.TotalAmount, &order.SellerPayout,
			&order.Status,
			&order.ShippingAddressID,
			&order.TrackingNumber, &order.Carrier,
			&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
			&order.CompletedAt, &order.CancelledAt,
			&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
			&order.CancellationReason,
			&order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID int64, newStatus string) error {
	var timestampCol string

	switch newStatus {
	case model.StatusPaid:
		timestampCol = ", payment_at = NOW()"
	case model.StatusShipped:
		timestampCol = ", shipped_at = NOW()"
	case model.StatusDelivered:
		timestampCol = ", delivered_at = NOW()"
	case model.StatusCompleted:
		timestampCol = ", completed_at = NOW()"
	case model.StatusCancelled, model.StatusRefunded:
		timestampCol = ", cancelled_at = NOW()"
	}

	query := fmt.Sprintf(`
		UPDATE orders
		SET status = $1%s
		WHERE id = $2
	`, timestampCol)

	result, err := r.db.Exec(ctx, query, newStatus, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// AddTrackingNumber adds tracking information to an order
func (r *OrderRepository) AddTrackingNumber(ctx context.Context, orderID int64, trackingNumber, carrier string) error {
	query := `
		UPDATE orders
		SET 
			tracking_number = $1,
			carrier = $2,
			status = $3,
			shipped_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.Exec(ctx, query, trackingNumber, carrier, model.StatusShipped, orderID)
	if err != nil {
		return fmt.Errorf("failed to add tracking number: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(ctx context.Context, orderID int64, reason string) error {
	query := `
		UPDATE orders
		SET 
			status = $1,
			cancellation_reason = $2,
			cancelled_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query, model.StatusCancelled, reason, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// GetOrderStatusHistory retrieves the status history for an order
func (r *OrderRepository) GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*model.OrderStatusHistory, error) {
	query := `
		SELECT 
			id, order_id, from_status, to_status,
			note, created_by, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status history: %w", err)
	}
	defer rows.Close()

	history := make([]*model.OrderStatusHistory, 0)
	for rows.Next() {
		h := &model.OrderStatusHistory{}
		err := rows.Scan(
			&h.ID, &h.OrderID, &h.FromStatus, &h.ToStatus,
			&h.Note, &h.CreatedBy, &h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

// GetOrderByMatchID retrieves an order by match ID
func (r *OrderRepository) GetOrderByMatchID(ctx context.Context, matchID int64) (*model.Order, error) {
	query := `
		SELECT 
			id, order_number, match_id,
			buyer_id, seller_id,
			product_id, size_id,
			price, quantity,
			buyer_fee, seller_fee, platform_fee,
			total_amount, seller_payout,
			status,
			shipping_address_id,
			tracking_number, carrier,
			payment_at, shipped_at, delivered_at, completed_at, cancelled_at,
			buyer_notes, seller_notes, admin_notes, cancellation_reason,
			created_at, updated_at
		FROM orders
		WHERE match_id = $1
	`

	order := &model.Order{}
	err := r.db.QueryRow(ctx, query, matchID).Scan(
		&order.ID, &order.OrderNumber, &order.MatchID,
		&order.BuyerID, &order.SellerID,
		&order.ProductID, &order.SizeID,
		&order.Price, &order.Quantity,
		&order.BuyerFee, &order.SellerFee, &order.PlatformFee,
		&order.TotalAmount, &order.SellerPayout,
		&order.Status,
		&order.ShippingAddressID,
		&order.TrackingNumber, &order.Carrier,
		&order.PaymentAt, &order.ShippedAt, &order.DeliveredAt,
		&order.CompletedAt, &order.CancelledAt,
		&order.BuyerNotes, &order.SellerNotes, &order.AdminNotes,
		&order.CancellationReason,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order doesn't exist yet for this match
		}
		return nil, fmt.Errorf("failed to get order by match ID: %w", err)
	}

	return order, nil
}
