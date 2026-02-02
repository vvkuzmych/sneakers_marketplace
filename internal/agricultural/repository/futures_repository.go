package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/agricultural/model"
)

type FuturesContractRepository struct {
	db *pgxpool.Pool
}

func NewFuturesContractRepository(db *pgxpool.Pool) *FuturesContractRepository {
	return &FuturesContractRepository{db: db}
}

// CreateContract creates a new futures contract
func (r *FuturesContractRepository) CreateContract(ctx context.Context, req *model.CreateFuturesContractRequest) (*model.FuturesContract, error) {
	// Generate contract number
	contractNumber := fmt.Sprintf("FUT-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	query := `
		INSERT INTO futures_contracts (
			product_id, contract_number, contract_type, strike_price, quantity, unit_of_measure,
			delivery_date, expiration_date, buyer_id, seller_id, status, margin_requirement, delivery_location
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active', $11, $12)
		RETURNING id, product_id, contract_number, contract_type, strike_price, quantity, unit_of_measure,
			contract_date, delivery_date, expiration_date, buyer_id, seller_id, status, settled_price, settled_at,
			margin_requirement, margin_posted, delivery_location, delivery_status, created_at, updated_at
	`

	contract := &model.FuturesContract{}
	err := r.db.QueryRow(ctx, query,
		req.ProductID, contractNumber, req.ContractType, req.StrikePrice, req.Quantity, req.UnitOfMeasure,
		req.DeliveryDate, req.ExpirationDate, req.BuyerID, req.SellerID, req.MarginRequirement, req.DeliveryLocation,
	).Scan(
		&contract.ID, &contract.ProductID, &contract.ContractNumber, &contract.ContractType, &contract.StrikePrice,
		&contract.Quantity, &contract.UnitOfMeasure, &contract.ContractDate, &contract.DeliveryDate, &contract.ExpirationDate,
		&contract.BuyerID, &contract.SellerID, &contract.Status, &contract.SettledPrice, &contract.SettledAt,
		&contract.MarginRequirement, &contract.MarginPosted, &contract.DeliveryLocation, &contract.DeliveryStatus,
		&contract.CreatedAt, &contract.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create futures contract: %w", err)
	}

	return contract, nil
}

// GetContractByID retrieves a contract by ID
func (r *FuturesContractRepository) GetContractByID(ctx context.Context, id int64) (*model.FuturesContract, error) {
	query := `
		SELECT id, product_id, contract_number, contract_type, strike_price, quantity, unit_of_measure,
			contract_date, delivery_date, expiration_date, buyer_id, seller_id, status, settled_price, settled_at,
			margin_requirement, margin_posted, delivery_location, delivery_status, created_at, updated_at
		FROM futures_contracts
		WHERE id = $1
	`

	contract := &model.FuturesContract{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&contract.ID, &contract.ProductID, &contract.ContractNumber, &contract.ContractType, &contract.StrikePrice,
		&contract.Quantity, &contract.UnitOfMeasure, &contract.ContractDate, &contract.DeliveryDate, &contract.ExpirationDate,
		&contract.BuyerID, &contract.SellerID, &contract.Status, &contract.SettledPrice, &contract.SettledAt,
		&contract.MarginRequirement, &contract.MarginPosted, &contract.DeliveryLocation, &contract.DeliveryStatus,
		&contract.CreatedAt, &contract.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return contract, nil
}

// ListContractsByUser retrieves all contracts for a user (buyer or seller)
func (r *FuturesContractRepository) ListContractsByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.FuturesContract, error) {
	query := `
		SELECT id, product_id, contract_number, contract_type, strike_price, quantity, unit_of_measure,
			contract_date, delivery_date, expiration_date, buyer_id, seller_id, status, settled_price, settled_at,
			margin_requirement, margin_posted, delivery_location, delivery_status, created_at, updated_at
		FROM futures_contracts
		WHERE buyer_id = $1 OR seller_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*model.FuturesContract
	for rows.Next() {
		contract := &model.FuturesContract{}
		err := rows.Scan(
			&contract.ID, &contract.ProductID, &contract.ContractNumber, &contract.ContractType, &contract.StrikePrice,
			&contract.Quantity, &contract.UnitOfMeasure, &contract.ContractDate, &contract.DeliveryDate, &contract.ExpirationDate,
			&contract.BuyerID, &contract.SellerID, &contract.Status, &contract.SettledPrice, &contract.SettledAt,
			&contract.MarginRequirement, &contract.MarginPosted, &contract.DeliveryLocation, &contract.DeliveryStatus,
			&contract.CreatedAt, &contract.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

// SettleContract settles a futures contract
func (r *FuturesContractRepository) SettleContract(ctx context.Context, id int64, settledPrice float64) error {
	query := `
		UPDATE futures_contracts
		SET status = 'settled',
			settled_price = $1,
			settled_at = NOW(),
			updated_at = NOW()
		WHERE id = $2 AND status = 'active'
	`

	result, err := r.db.Exec(ctx, query, settledPrice, id)
	if err != nil {
		return fmt.Errorf("failed to settle contract: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("contract not found or already settled")
	}

	return nil
}

// CancelContract cancels a futures contract
func (r *FuturesContractRepository) CancelContract(ctx context.Context, id int64) error {
	query := `
		UPDATE futures_contracts
		SET status = 'cancelled',
			updated_at = NOW()
		WHERE id = $1 AND status = 'active'
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel contract: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("contract not found or cannot be cancelled")
	}

	return nil
}

// ExpireContracts marks expired contracts
func (r *FuturesContractRepository) ExpireContracts(ctx context.Context) (int, error) {
	query := `
		UPDATE futures_contracts
		SET status = 'expired',
			updated_at = NOW()
		WHERE status = 'active' AND expiration_date < NOW()
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to expire contracts: %w", err)
	}

	return int(result.RowsAffected()), nil
}
