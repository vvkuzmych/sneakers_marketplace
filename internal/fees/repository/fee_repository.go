package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/fees/model"
)

type FeeRepository struct {
	db *pgxpool.Pool
}

func NewFeeRepository(db *pgxpool.Pool) *FeeRepository {
	return &FeeRepository{db: db}
}

// GetFeeConfig retrieves fee configuration for a specific vertical
func (r *FeeRepository) GetFeeConfig(ctx context.Context, vertical string) (*model.FeeConfig, error) {
	query := `
		SELECT 
			id, vertical, transaction_fee_percent, processing_fee_fixed,
			authentication_fee, shipping_buyer_charge, shipping_seller_cost,
			min_transaction_fee, max_transaction_fee, created_at, updated_at
		FROM fee_configs
		WHERE vertical = $1
	`

	var config model.FeeConfig
	err := r.db.QueryRow(ctx, query, vertical).Scan(
		&config.ID,
		&config.Vertical,
		&config.TransactionFeePercent,
		&config.ProcessingFeeFixed,
		&config.AuthenticationFee,
		&config.ShippingBuyerCharge,
		&config.ShippingSellerCost,
		&config.MinTransactionFee,
		&config.MaxTransactionFee,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get fee config for vertical %s: %w", vertical, err)
	}

	return &config, nil
}

// GetAllFeeConfigs retrieves all fee configurations
func (r *FeeRepository) GetAllFeeConfigs(ctx context.Context) ([]*model.FeeConfig, error) {
	query := `
		SELECT 
			id, vertical, transaction_fee_percent, processing_fee_fixed,
			authentication_fee, shipping_buyer_charge, shipping_seller_cost,
			min_transaction_fee, max_transaction_fee, created_at, updated_at
		FROM fee_configs
		ORDER BY vertical
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee configs: %w", err)
	}
	defer rows.Close()

	var configs []*model.FeeConfig
	for rows.Next() {
		var config model.FeeConfig
		err := rows.Scan(
			&config.ID,
			&config.Vertical,
			&config.TransactionFeePercent,
			&config.ProcessingFeeFixed,
			&config.AuthenticationFee,
			&config.ShippingBuyerCharge,
			&config.ShippingSellerCost,
			&config.MinTransactionFee,
			&config.MaxTransactionFee,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fee config: %w", err)
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// CreateTransactionFee records a transaction fee to the database
func (r *FeeRepository) CreateTransactionFee(ctx context.Context, fee *model.TransactionFee) error {
	query := `
		INSERT INTO transaction_fees (
			match_id, order_id, vertical, sale_price,
			buyer_processing_fee, buyer_shipping_fee, buyer_total,
			seller_transaction_fee, seller_authentication_fee, seller_shipping_cost, seller_payout,
			platform_revenue, fee_config_snapshot
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	snapshotJSON, err := json.Marshal(fee.FeeConfigSnapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal fee config snapshot: %w", err)
	}

	err = r.db.QueryRow(ctx, query,
		fee.MatchID,
		fee.OrderID,
		fee.Vertical,
		fee.SalePrice,
		fee.BuyerProcessingFee,
		fee.BuyerShippingFee,
		fee.BuyerTotal,
		fee.SellerTransactionFee,
		fee.SellerAuthenticationFee,
		fee.SellerShippingCost,
		fee.SellerPayout,
		fee.PlatformRevenue,
		snapshotJSON,
	).Scan(&fee.ID, &fee.CreatedAt, &fee.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transaction fee: %w", err)
	}

	return nil
}

// GetTransactionFeeByMatchID retrieves transaction fee by match ID
func (r *FeeRepository) GetTransactionFeeByMatchID(ctx context.Context, matchID int64) (*model.TransactionFee, error) {
	query := `
		SELECT 
			id, match_id, order_id, vertical, sale_price,
			buyer_processing_fee, buyer_shipping_fee, buyer_total,
			seller_transaction_fee, seller_authentication_fee, seller_shipping_cost, seller_payout,
			platform_revenue, fee_config_snapshot, created_at, updated_at
		FROM transaction_fees
		WHERE match_id = $1
	`

	var fee model.TransactionFee
	var snapshotJSON []byte

	err := r.db.QueryRow(ctx, query, matchID).Scan(
		&fee.ID,
		&fee.MatchID,
		&fee.OrderID,
		&fee.Vertical,
		&fee.SalePrice,
		&fee.BuyerProcessingFee,
		&fee.BuyerShippingFee,
		&fee.BuyerTotal,
		&fee.SellerTransactionFee,
		&fee.SellerAuthenticationFee,
		&fee.SellerShippingCost,
		&fee.SellerPayout,
		&fee.PlatformRevenue,
		&snapshotJSON,
		&fee.CreatedAt,
		&fee.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction fee for match %d: %w", matchID, err)
	}

	// Parse JSONB snapshot
	if len(snapshotJSON) > 0 {
		if err := json.Unmarshal(snapshotJSON, &fee.FeeConfigSnapshot); err != nil {
			return nil, fmt.Errorf("failed to unmarshal fee config snapshot: %w", err)
		}
	}

	return &fee, nil
}

// GetTotalRevenue calculates total platform revenue for a date range
func (r *FeeRepository) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(platform_revenue), 0)
		FROM transaction_fees
		WHERE created_at >= $1 AND created_at < $2
	`

	var total float64
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}

	return total, nil
}

// GetRevenueByVertical calculates revenue per vertical for a date range
func (r *FeeRepository) GetRevenueByVertical(ctx context.Context, startDate, endDate time.Time) (map[string]float64, error) {
	query := `
		SELECT vertical, COALESCE(SUM(platform_revenue), 0) as revenue
		FROM transaction_fees
		WHERE created_at >= $1 AND created_at < $2
		GROUP BY vertical
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by vertical: %w", err)
	}
	defer rows.Close()

	revenueMap := make(map[string]float64)
	for rows.Next() {
		var vertical string
		var revenue float64
		if err := rows.Scan(&vertical, &revenue); err != nil {
			return nil, fmt.Errorf("failed to scan revenue row: %w", err)
		}
		revenueMap[vertical] = revenue
	}

	return revenueMap, nil
}

// GetTransactionCount returns count of transactions for a date range
func (r *FeeRepository) GetTransactionCount(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM transaction_fees
		WHERE created_at >= $1 AND created_at < $2
	`

	var count int64
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}

	return count, nil
}
