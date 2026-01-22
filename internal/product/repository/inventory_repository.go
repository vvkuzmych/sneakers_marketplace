package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/product/model"
)

// InventoryRepository handles inventory operations
type InventoryRepository struct {
	db *pgxpool.Pool
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// AddSize adds a new size with initial quantity
func (r *InventoryRepository) AddSize(ctx context.Context, size *model.Size) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Insert size
	query := `
		INSERT INTO sizes (product_id, size, quantity, reserved)
		VALUES ($1, $2, $3, 0)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		size.ProductID,
		size.Size,
		size.Quantity,
	).Scan(&size.ID, &size.CreatedAt, &size.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to add size: %w", err)
	}

	// Log inventory transaction
	if err := r.logInventoryTransaction(ctx, tx, size.ID, "addition", size.Quantity, 0, size.Quantity, "", "Initial inventory"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetSizesByProductID retrieves all sizes for a product
func (r *InventoryRepository) GetSizesByProductID(ctx context.Context, productID int64) ([]*model.Size, error) {
	query := `
		SELECT id, product_id, size, quantity, reserved, created_at, updated_at
		FROM sizes
		WHERE product_id = $1
		ORDER BY size ASC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}
	defer rows.Close()

	var sizes []*model.Size
	for rows.Next() {
		size := &model.Size{}
		err := rows.Scan(
			&size.ID,
			&size.ProductID,
			&size.Size,
			&size.Quantity,
			&size.Reserved,
			&size.CreatedAt,
			&size.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan size: %w", err)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

// GetSizeByID retrieves a size by ID
func (r *InventoryRepository) GetSizeByID(ctx context.Context, sizeID int64) (*model.Size, error) {
	query := `
		SELECT id, product_id, size, quantity, reserved, created_at, updated_at
		FROM sizes
		WHERE id = $1
	`

	size := &model.Size{}
	err := r.db.QueryRow(ctx, query, sizeID).Scan(
		&size.ID,
		&size.ProductID,
		&size.Size,
		&size.Quantity,
		&size.Reserved,
		&size.CreatedAt,
		&size.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return size, nil
}

// UpdateInventory updates the quantity for a size
func (r *InventoryRepository) UpdateInventory(ctx context.Context, sizeID int64, newQuantity int, notes string) (*model.Size, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get current size
	var size model.Size
	query := `SELECT id, product_id, size, quantity, reserved, created_at, updated_at FROM sizes WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, sizeID).Scan(
		&size.ID,
		&size.ProductID,
		&size.Size,
		&size.Quantity,
		&size.Reserved,
		&size.CreatedAt,
		&size.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	oldQuantity := size.Quantity
	quantityChange := newQuantity - oldQuantity

	// Update quantity
	updateQuery := `
		UPDATE sizes
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, updateQuery, newQuantity, sizeID).Scan(&size.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	size.Quantity = newQuantity

	// Log transaction
	transactionType := "addition"
	if quantityChange < 0 {
		transactionType = "removal"
	}

	if err := r.logInventoryTransaction(ctx, tx, sizeID, transactionType, quantityChange, oldQuantity, newQuantity, "", notes); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &size, nil
}

// ReserveInventory reserves inventory for an order
func (r *InventoryRepository) ReserveInventory(ctx context.Context, sizeID int64, quantity int, orderID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get and lock size
	var size model.Size
	query := `SELECT id, quantity, reserved FROM sizes WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, sizeID).Scan(&size.ID, &size.Quantity, &size.Reserved)
	if err != nil {
		return fmt.Errorf("failed to get size: %w", err)
	}

	// Check if enough inventory available
	available := size.Quantity - size.Reserved
	if available < quantity {
		return fmt.Errorf("insufficient inventory: available=%d, requested=%d", available, quantity)
	}

	// Reserve inventory
	updateQuery := `
		UPDATE sizes
		SET reserved = reserved + $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = tx.Exec(ctx, updateQuery, quantity, sizeID)
	if err != nil {
		return fmt.Errorf("failed to reserve inventory: %w", err)
	}

	// Log transaction
	if err := r.logInventoryTransaction(ctx, tx, sizeID, "reservation", -quantity, size.Quantity, size.Quantity, orderID, fmt.Sprintf("Reserved for order %s", orderID)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ReleaseInventory releases reserved inventory
func (r *InventoryRepository) ReleaseInventory(ctx context.Context, sizeID int64, quantity int, orderID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get and lock size
	var size model.Size
	query := `SELECT id, quantity, reserved FROM sizes WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, sizeID).Scan(&size.ID, &size.Quantity, &size.Reserved)
	if err != nil {
		return fmt.Errorf("failed to get size: %w", err)
	}

	// Ensure we don't release more than reserved
	if size.Reserved < quantity {
		return fmt.Errorf("cannot release more than reserved: reserved=%d, release=%d", size.Reserved, quantity)
	}

	// Release inventory
	updateQuery := `
		UPDATE sizes
		SET reserved = reserved - $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = tx.Exec(ctx, updateQuery, quantity, sizeID)
	if err != nil {
		return fmt.Errorf("failed to release inventory: %w", err)
	}

	// Log transaction
	if err := r.logInventoryTransaction(ctx, tx, sizeID, "release", quantity, size.Quantity, size.Quantity, orderID, fmt.Sprintf("Released from order %s", orderID)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CompleteSale completes a sale by reducing both quantity and reserved
func (r *InventoryRepository) CompleteSale(ctx context.Context, sizeID int64, quantity int, orderID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get and lock size
	var size model.Size
	query := `SELECT id, quantity, reserved FROM sizes WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, sizeID).Scan(&size.ID, &size.Quantity, &size.Reserved)
	if err != nil {
		return fmt.Errorf("failed to get size: %w", err)
	}

	oldQuantity := size.Quantity

	// Reduce both quantity and reserved
	updateQuery := `
		UPDATE sizes
		SET quantity = quantity - $1, reserved = reserved - $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = tx.Exec(ctx, updateQuery, quantity, sizeID)
	if err != nil {
		return fmt.Errorf("failed to complete sale: %w", err)
	}

	// Log transaction
	if err := r.logInventoryTransaction(ctx, tx, sizeID, "sale", -quantity, oldQuantity, oldQuantity-quantity, orderID, fmt.Sprintf("Sale completed for order %s", orderID)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// logInventoryTransaction logs an inventory change
func (r *InventoryRepository) logInventoryTransaction(ctx context.Context, tx pgx.Tx, sizeID int64, transactionType string, quantityChange, quantityBefore, quantityAfter int, referenceID, notes string) error {
	query := `
		INSERT INTO inventory_transactions (size_id, transaction_type, quantity_change, quantity_before, quantity_after, reference_id, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(ctx, query,
		sizeID,
		transactionType,
		quantityChange,
		quantityBefore,
		quantityAfter,
		referenceID,
		notes,
	)

	if err != nil {
		return fmt.Errorf("failed to log inventory transaction: %w", err)
	}

	return nil
}

// GetInventoryTransactions retrieves transaction history for a size
func (r *InventoryRepository) GetInventoryTransactions(ctx context.Context, sizeID int64, limit int) ([]*model.InventoryTransaction, error) {
	query := `
		SELECT id, size_id, transaction_type, quantity_change, quantity_before, quantity_after, reference_id, notes, created_at
		FROM inventory_transactions
		WHERE size_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, sizeID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.InventoryTransaction
	for rows.Next() {
		txn := &model.InventoryTransaction{}
		err := rows.Scan(
			&txn.ID,
			&txn.SizeID,
			&txn.TransactionType,
			&txn.QuantityChange,
			&txn.QuantityBefore,
			&txn.QuantityAfter,
			&txn.ReferenceID,
			&txn.Notes,
			&txn.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}
