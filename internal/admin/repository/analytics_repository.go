package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/model"
)

// ==================== Platform Statistics ====================

// GetPlatformStats retrieves overall platform statistics
func (r *AdminRepository) GetPlatformStats(ctx context.Context) (*model.PlatformStats, error) {
	query := `
		WITH today AS (
			SELECT CURRENT_DATE as date
		)
		SELECT 
			(SELECT COUNT(*) FROM users WHERE is_active = true) as total_users,
			(SELECT COUNT(DISTINCT user_id) FROM sessions WHERE created_at >= (SELECT date FROM today)) as active_users_today,
			(SELECT COUNT(*) FROM products) as total_products,
			(SELECT COUNT(*) FROM products WHERE is_active = true) as active_products,
			(SELECT COUNT(*) FROM orders) as total_orders,
			(SELECT COUNT(*) FROM orders WHERE DATE(created_at) = (SELECT date FROM today)) as orders_today,
			(SELECT COALESCE(SUM(total), 0) FROM orders WHERE status IN ('delivered', 'completed')) as total_revenue,
			(SELECT COALESCE(SUM(total), 0) FROM orders WHERE status IN ('delivered', 'completed') AND DATE(created_at) = (SELECT date FROM today)) as revenue_today,
			(SELECT COALESCE(SUM(buyer_fee + seller_fee), 0) FROM orders WHERE status IN ('delivered', 'completed')) as total_fees,
			(SELECT COUNT(*) FROM matches) as total_matches,
			(SELECT COUNT(*) FROM matches WHERE DATE(created_at) = (SELECT date FROM today)) as matches_today
	`

	var stats model.PlatformStats
	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalUsers, &stats.ActiveUsersToday, &stats.TotalProducts, &stats.ActiveProducts,
		&stats.TotalOrders, &stats.OrdersToday, &stats.TotalRevenue, &stats.RevenueToday,
		&stats.TotalFeesCollected, &stats.TotalMatches, &stats.MatchesToday,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform stats: %w", err)
	}

	return &stats, nil
}

// ==================== Revenue Report ====================

// GetRevenueReport retrieves revenue report grouped by time period
func (r *AdminRepository) GetRevenueReport(ctx context.Context, params model.GetRevenueReportParams) (*model.RevenueReport, error) {
	// Determine date truncation based on groupBy
	dateTrunc := "day"
	switch params.GroupBy {
	case "week":
		dateTrunc = "week"
	case "month":
		dateTrunc = "month"
	default:
		dateTrunc = "day"
	}

	query := fmt.Sprintf(`
		SELECT 
			DATE_TRUNC('%s', created_at)::date as period,
			COALESCE(SUM(total), 0) as revenue,
			COALESCE(SUM(buyer_fee + seller_fee), 0) as fees,
			COUNT(*) as order_count
		FROM orders
		WHERE created_at >= $1 AND created_at <= $2
		  AND status IN ('delivered', 'completed')
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY period ASC
	`, dateTrunc, dateTrunc)

	rows, err := r.db.Query(ctx, query, params.DateFrom, params.DateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue report: %w", err)
	}
	defer rows.Close()

	dataPoints := []model.RevenueDataPoint{}
	var totalRevenue, totalFees float64

	for rows.Next() {
		var dp model.RevenueDataPoint
		var period time.Time
		err := rows.Scan(&period, &dp.Revenue, &dp.Fees, &dp.OrderCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan revenue data: %w", err)
		}

		// Format label based on groupBy
		switch params.GroupBy {
		case "month":
			dp.Label = period.Format("2006-01")
		case "week":
			dp.Label = fmt.Sprintf("Week of %s", period.Format("2006-01-02"))
		default:
			dp.Label = period.Format("2006-01-02")
		}

		totalRevenue += dp.Revenue
		totalFees += dp.Fees
		dataPoints = append(dataPoints, dp)
	}

	return &model.RevenueReport{
		DataPoints:   dataPoints,
		TotalRevenue: totalRevenue,
		TotalFees:    totalFees,
		DateFrom:     params.DateFrom,
		DateTo:       params.DateTo,
		GroupBy:      params.GroupBy,
	}, nil
}

// ==================== User Activity Report ====================

// GetUserActivityReport retrieves user activity statistics
func (r *AdminRepository) GetUserActivityReport(ctx context.Context, params model.GetUserActivityReportParams) (*model.UserActivityReport, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM users WHERE created_at >= $1 AND created_at <= $2) as new_users,
			(SELECT COUNT(DISTINCT user_id) FROM sessions WHERE created_at >= $1 AND created_at <= $2) as active_users,
			(SELECT COUNT(*) FROM bids WHERE created_at >= $1 AND created_at <= $2) as total_bids,
			(SELECT COUNT(*) FROM asks WHERE created_at >= $1 AND created_at <= $2) as total_asks,
			(SELECT COUNT(*) FROM matches WHERE created_at >= $1 AND created_at <= $2) as total_matches
	`

	var report model.UserActivityReport
	report.DateFrom = params.DateFrom
	report.DateTo = params.DateTo

	err := r.db.QueryRow(ctx, query, params.DateFrom, params.DateTo).Scan(
		&report.NewUsers, &report.ActiveUsers, &report.TotalBidsPlaced,
		&report.TotalAsksPlaced, &report.TotalMatchesCreated,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity report: %w", err)
	}

	return &report, nil
}

// ==================== Order Management ====================

// ListAllOrders retrieves all orders for admin view
func (r *AdminRepository) ListAllOrders(ctx context.Context, params model.ListOrdersParams) ([]model.OrderSummary, int32, error) {
	query := `
		SELECT 
			o.id, o.order_number, o.buyer_id, o.seller_id,
			buyer.email as buyer_email, seller.email as seller_email,
			o.product_id, p.name as product_name,
			o.subtotal, o.buyer_fee, o.seller_fee, o.total, o.status, o.created_at
		FROM orders o
		JOIN users buyer ON buyer.id = o.buyer_id
		JOIN users seller ON seller.id = o.seller_id
		JOIN products p ON p.id = o.product_id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Apply filters
	if params.Status != "" && params.Status != "all" {
		query += fmt.Sprintf(" AND o.status = $%d", argCount)
		args = append(args, params.Status)
		argCount++
	}

	if params.DateFrom != nil {
		query += fmt.Sprintf(" AND o.created_at >= $%d", argCount)
		args = append(args, *params.DateFrom)
		argCount++
	}

	if params.DateTo != nil {
		query += fmt.Sprintf(" AND o.created_at <= $%d", argCount)
		args = append(args, *params.DateTo)
		argCount++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS filtered"
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Apply sorting
	sortBy := "o.created_at"
	if params.SortBy == "total_amount" {
		sortBy = "o.total"
	}
	sortOrder := "DESC"
	if params.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	offset := (params.Page - 1) * params.PageSize
	args = append(args, params.PageSize, offset)

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := []model.OrderSummary{}
	for rows.Next() {
		var o model.OrderSummary
		err := rows.Scan(
			&o.ID, &o.OrderNumber, &o.BuyerID, &o.SellerID,
			&o.BuyerEmail, &o.SellerEmail, &o.ProductID, &o.ProductName,
			&o.Subtotal, &o.BuyerFee, &o.SellerFee, &o.Total, &o.Status, &o.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// GetOrderDetails retrieves detailed order information
func (r *AdminRepository) GetOrderDetails(ctx context.Context, orderID int64) (*model.OrderSummary, []model.OrderStatusChange, error) {
	// Get order
	orderQuery := `
		SELECT 
			o.id, o.order_number, o.buyer_id, o.seller_id,
			buyer.email as buyer_email, seller.email as seller_email,
			o.product_id, p.name as product_name,
			o.subtotal, o.buyer_fee, o.seller_fee, o.total, o.status, o.created_at
		FROM orders o
		JOIN users buyer ON buyer.id = o.buyer_id
		JOIN users seller ON seller.id = o.seller_id
		JOIN products p ON p.id = o.product_id
		WHERE o.id = $1
	`

	var order model.OrderSummary
	err := r.db.QueryRow(ctx, orderQuery, orderID).Scan(
		&order.ID, &order.OrderNumber, &order.BuyerID, &order.SellerID,
		&order.BuyerEmail, &order.SellerEmail, &order.ProductID, &order.ProductName,
		&order.Subtotal, &order.BuyerFee, &order.SellerFee, &order.Total, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get status history
	historyQuery := `
		SELECT status, notes, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, historyQuery, orderID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get status history: %w", err)
	}
	defer rows.Close()

	history := []model.OrderStatusChange{}
	for rows.Next() {
		var change model.OrderStatusChange
		change.ChangedBy = "system" // Default
		err := rows.Scan(&change.Status, &change.Notes, &change.ChangedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan status change: %w", err)
		}
		history = append(history, change)
	}

	return &order, history, nil
}

// CancelOrder cancels an order
func (r *AdminRepository) CancelOrder(ctx context.Context, orderID int64, reason string) error {
	query := `
		UPDATE orders 
		SET status = 'cancelled',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status NOT IN ('cancelled', 'delivered', 'completed')
	`

	result, err := r.db.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found or cannot be canceled")
	}

	// Add status history
	historyQuery := `
		INSERT INTO order_status_history (order_id, status, notes, created_at)
		VALUES ($1, 'cancelled', $2, CURRENT_TIMESTAMP)
	`
	_, err = r.db.Exec(ctx, historyQuery, orderID, reason)
	if err != nil {
		return fmt.Errorf("failed to add status history: %w", err)
	}

	return nil
}

// ==================== Product Management ====================

// ListAllProducts retrieves all products for admin view
func (r *AdminRepository) ListAllProducts(ctx context.Context, params model.ListProductsParams) ([]model.ProductSummary, int32, error) {
	query := `
		SELECT 
			p.id, p.sku, p.name, p.brand, p.model, p.retail_price, 
			p.is_active, p.is_featured, p.created_at,
			COALESCE(COUNT(DISTINCT b.id), 0) as total_bids,
			COALESCE(COUNT(DISTINCT a.id), 0) as total_asks,
			COALESCE(COUNT(DISTINCT m.id), 0) as total_matches,
			COALESCE(MAX(b.price), 0) as highest_bid,
			COALESCE(MIN(a.price), 0) as lowest_ask
		FROM products p
		LEFT JOIN bids b ON b.product_id = p.id AND b.status = 'active'
		LEFT JOIN asks a ON a.product_id = p.id AND a.status = 'active'
		LEFT JOIN matches m ON m.product_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Apply filters
	if params.Status == "active" {
		query += " AND p.is_active = true"
	} else if params.Status == "hidden" {
		query += " AND p.is_active = false"
	} else if params.Status == "featured" {
		query += " AND p.is_featured = true"
	}

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.brand ILIKE $%d OR p.sku ILIKE $%d)",
			argCount, argCount, argCount)
		args = append(args, searchPattern)
		argCount++
	}

	query += " GROUP BY p.id"

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS filtered"
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	offset := (params.Page - 1) * params.PageSize
	args = append(args, params.PageSize, offset)

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []model.ProductSummary{}
	for rows.Next() {
		var p model.ProductSummary
		err := rows.Scan(
			&p.ID, &p.SKU, &p.Name, &p.Brand, &p.Model, &p.RetailPrice,
			&p.IsActive, &p.IsFeatured, &p.CreatedAt,
			&p.TotalBids, &p.TotalAsks, &p.TotalMatches, &p.HighestBid, &p.LowestAsk,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}

// FeatureProduct marks a product as featured
func (r *AdminRepository) FeatureProduct(ctx context.Context, productID int64) error {
	query := `UPDATE products SET is_featured = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to feature product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// HideProduct hides a product
func (r *AdminRepository) HideProduct(ctx context.Context, productID int64) error {
	query := `UPDATE products SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to hide product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
