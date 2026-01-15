package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/model"
)

// BiddingRepository handles database operations for bidding
type BiddingRepository struct {
	db *pgxpool.Pool
}

// NewBiddingRepository creates a new bidding repository
func NewBiddingRepository(db *pgxpool.Pool) *BiddingRepository {
	return &BiddingRepository{db: db}
}

// PlaceBid creates a new bid
func (r *BiddingRepository) PlaceBid(ctx context.Context, bid *model.Bid) error {
	query := `
		INSERT INTO bids (user_id, product_id, size_id, price, quantity, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		bid.UserID,
		bid.ProductID,
		bid.SizeID,
		bid.Price,
		bid.Quantity,
		bid.Status,
		bid.ExpiresAt,
	).Scan(&bid.ID, &bid.CreatedAt, &bid.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to place bid: %w", err)
	}

	return nil
}

// GetBidByID retrieves a bid by ID
func (r *BiddingRepository) GetBidByID(ctx context.Context, bidID int64) (*model.Bid, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM bids
		WHERE id = $1
	`

	bid := &model.Bid{}
	err := r.db.QueryRow(ctx, query, bidID).Scan(
		&bid.ID,
		&bid.UserID,
		&bid.ProductID,
		&bid.SizeID,
		&bid.Price,
		&bid.Quantity,
		&bid.Status,
		&bid.ExpiresAt,
		&bid.MatchedAt,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get bid: %w", err)
	}

	return bid, nil
}

// GetUserBids retrieves bids for a user
func (r *BiddingRepository) GetUserBids(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Bid, int64, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM bids
		WHERE user_id = $1
	`
	countQuery := `SELECT COUNT(*) FROM bids WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bids: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user bids: %w", err)
	}
	defer rows.Close()

	var bids []*model.Bid
	for rows.Next() {
		bid := &model.Bid{}
		err := rows.Scan(
			&bid.ID,
			&bid.UserID,
			&bid.ProductID,
			&bid.SizeID,
			&bid.Price,
			&bid.Quantity,
			&bid.Status,
			&bid.ExpiresAt,
			&bid.MatchedAt,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	return bids, total, nil
}

// GetProductBids retrieves bids for a product/size
func (r *BiddingRepository) GetProductBids(ctx context.Context, productID, sizeID int64, status string, page, pageSize int) ([]*model.Bid, int64, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM bids
		WHERE product_id = $1 AND size_id = $2
	`
	countQuery := `SELECT COUNT(*) FROM bids WHERE product_id = $1 AND size_id = $2`

	args := []interface{}{productID, sizeID}
	argIdx := 3

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bids: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY price DESC, created_at ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product bids: %w", err)
	}
	defer rows.Close()

	var bids []*model.Bid
	for rows.Next() {
		bid := &model.Bid{}
		err := rows.Scan(
			&bid.ID,
			&bid.UserID,
			&bid.ProductID,
			&bid.SizeID,
			&bid.Price,
			&bid.Quantity,
			&bid.Status,
			&bid.ExpiresAt,
			&bid.MatchedAt,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	return bids, total, nil
}

// UpdateBidStatus updates bid status
func (r *BiddingRepository) UpdateBidStatus(ctx context.Context, tx pgx.Tx, bidID int64, status string) error {
	query := `UPDATE bids SET status = $1, updated_at = NOW() WHERE id = $2`

	var result pgconn.CommandTag
	var err error

	if tx != nil {
		result, err = tx.Exec(ctx, query, status, bidID)
	} else {
		result, err = r.db.Exec(ctx, query, status, bidID)
	}

	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("bid not found")
	}

	return nil
}

// GetHighestBid retrieves the highest active bid for a product/size
func (r *BiddingRepository) GetHighestBid(ctx context.Context, productID, sizeID int64) (*model.Bid, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM bids
		WHERE product_id = $1 AND size_id = $2 AND status = 'active'
		ORDER BY price DESC, created_at ASC
		LIMIT 1
	`

	bid := &model.Bid{}
	err := r.db.QueryRow(ctx, query, productID, sizeID).Scan(
		&bid.ID,
		&bid.UserID,
		&bid.ProductID,
		&bid.SizeID,
		&bid.Price,
		&bid.Quantity,
		&bid.Status,
		&bid.ExpiresAt,
		&bid.MatchedAt,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // No bids found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get highest bid: %w", err)
	}

	return bid, nil
}

// PlaceAsk creates a new ask
func (r *BiddingRepository) PlaceAsk(ctx context.Context, ask *model.Ask) error {
	query := `
		INSERT INTO asks (user_id, product_id, size_id, price, quantity, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		ask.UserID,
		ask.ProductID,
		ask.SizeID,
		ask.Price,
		ask.Quantity,
		ask.Status,
		ask.ExpiresAt,
	).Scan(&ask.ID, &ask.CreatedAt, &ask.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to place ask: %w", err)
	}

	return nil
}

// GetAskByID retrieves an ask by ID
func (r *BiddingRepository) GetAskByID(ctx context.Context, askID int64) (*model.Ask, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM asks
		WHERE id = $1
	`

	ask := &model.Ask{}
	err := r.db.QueryRow(ctx, query, askID).Scan(
		&ask.ID,
		&ask.UserID,
		&ask.ProductID,
		&ask.SizeID,
		&ask.Price,
		&ask.Quantity,
		&ask.Status,
		&ask.ExpiresAt,
		&ask.MatchedAt,
		&ask.CreatedAt,
		&ask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get ask: %w", err)
	}

	return ask, nil
}

// GetUserAsks retrieves asks for a user
func (r *BiddingRepository) GetUserAsks(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Ask, int64, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM asks
		WHERE user_id = $1
	`
	countQuery := `SELECT COUNT(*) FROM asks WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count asks: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user asks: %w", err)
	}
	defer rows.Close()

	var asks []*model.Ask
	for rows.Next() {
		ask := &model.Ask{}
		err := rows.Scan(
			&ask.ID,
			&ask.UserID,
			&ask.ProductID,
			&ask.SizeID,
			&ask.Price,
			&ask.Quantity,
			&ask.Status,
			&ask.ExpiresAt,
			&ask.MatchedAt,
			&ask.CreatedAt,
			&ask.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan ask: %w", err)
		}
		asks = append(asks, ask)
	}

	return asks, total, nil
}

// GetProductAsks retrieves asks for a product/size
func (r *BiddingRepository) GetProductAsks(ctx context.Context, productID, sizeID int64, status string, page, pageSize int) ([]*model.Ask, int64, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM asks
		WHERE product_id = $1 AND size_id = $2
	`
	countQuery := `SELECT COUNT(*) FROM asks WHERE product_id = $1 AND size_id = $2`

	args := []interface{}{productID, sizeID}
	argIdx := 3

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count asks: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY price ASC, created_at ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product asks: %w", err)
	}
	defer rows.Close()

	var asks []*model.Ask
	for rows.Next() {
		ask := &model.Ask{}
		err := rows.Scan(
			&ask.ID,
			&ask.UserID,
			&ask.ProductID,
			&ask.SizeID,
			&ask.Price,
			&ask.Quantity,
			&ask.Status,
			&ask.ExpiresAt,
			&ask.MatchedAt,
			&ask.CreatedAt,
			&ask.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan ask: %w", err)
		}
		asks = append(asks, ask)
	}

	return asks, total, nil
}

// UpdateAskStatus updates ask status
func (r *BiddingRepository) UpdateAskStatus(ctx context.Context, tx pgx.Tx, askID int64, status string) error {
	query := `UPDATE asks SET status = $1, updated_at = NOW() WHERE id = $2`

	var result pgconn.CommandTag
	var err error

	if tx != nil {
		result, err = tx.Exec(ctx, query, status, askID)
	} else {
		result, err = r.db.Exec(ctx, query, status, askID)
	}

	if err != nil {
		return fmt.Errorf("failed to update ask status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ask not found")
	}

	return nil
}

// GetLowestAsk retrieves the lowest active ask for a product/size
func (r *BiddingRepository) GetLowestAsk(ctx context.Context, productID, sizeID int64) (*model.Ask, error) {
	query := `
		SELECT id, user_id, product_id, size_id, price, quantity, status, 
		       expires_at, matched_at, created_at, updated_at
		FROM asks
		WHERE product_id = $1 AND size_id = $2 AND status = 'active'
		ORDER BY price ASC, created_at ASC
		LIMIT 1
	`

	ask := &model.Ask{}
	err := r.db.QueryRow(ctx, query, productID, sizeID).Scan(
		&ask.ID,
		&ask.UserID,
		&ask.ProductID,
		&ask.SizeID,
		&ask.Price,
		&ask.Quantity,
		&ask.Status,
		&ask.ExpiresAt,
		&ask.MatchedAt,
		&ask.CreatedAt,
		&ask.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // No asks found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get lowest ask: %w", err)
	}

	return ask, nil
}

// CreateMatch creates a new match
func (r *BiddingRepository) CreateMatch(ctx context.Context, tx pgx.Tx, match *model.Match) error {
	query := `
		INSERT INTO matches (bid_id, ask_id, buyer_id, seller_id, product_id, size_id, price, quantity, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	err := tx.QueryRow(ctx, query,
		match.BidID,
		match.AskID,
		match.BuyerID,
		match.SellerID,
		match.ProductID,
		match.SizeID,
		match.Price,
		match.Quantity,
		match.Status,
	).Scan(&match.ID, &match.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}

	return nil
}

// GetMatchByID retrieves a match by ID
func (r *BiddingRepository) GetMatchByID(ctx context.Context, matchID int64) (*model.Match, error) {
	query := `
		SELECT id, bid_id, ask_id, buyer_id, seller_id, product_id, size_id, 
		       price, quantity, status, completed_at, created_at
		FROM matches
		WHERE id = $1
	`

	match := &model.Match{}
	err := r.db.QueryRow(ctx, query, matchID).Scan(
		&match.ID,
		&match.BidID,
		&match.AskID,
		&match.BuyerID,
		&match.SellerID,
		&match.ProductID,
		&match.SizeID,
		&match.Price,
		&match.Quantity,
		&match.Status,
		&match.CompletedAt,
		&match.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return match, nil
}

// GetUserMatches retrieves matches for a user
func (r *BiddingRepository) GetUserMatches(ctx context.Context, userID int64, asBuyer, asSeller bool, page, pageSize int) ([]*model.Match, int64, error) {
	query := `
		SELECT id, bid_id, ask_id, buyer_id, seller_id, product_id, size_id, 
		       price, quantity, status, completed_at, created_at
		FROM matches
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM matches WHERE 1=1`

	args := []interface{}{}
	argIdx := 1

	if asBuyer && !asSeller {
		query += fmt.Sprintf(" AND buyer_id = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND buyer_id = $%d", argIdx)
		args = append(args, userID)
		argIdx++
	} else if asSeller && !asBuyer {
		query += fmt.Sprintf(" AND seller_id = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND seller_id = $%d", argIdx)
		args = append(args, userID)
		argIdx++
	} else {
		query += fmt.Sprintf(" AND (buyer_id = $%d OR seller_id = $%d)", argIdx, argIdx)
		countQuery += fmt.Sprintf(" AND (buyer_id = $%d OR seller_id = $%d)", argIdx, argIdx)
		args = append(args, userID)
		argIdx++
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count matches: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user matches: %w", err)
	}
	defer rows.Close()

	var matches []*model.Match
	for rows.Next() {
		match := &model.Match{}
		err := rows.Scan(
			&match.ID,
			&match.BidID,
			&match.AskID,
			&match.BuyerID,
			&match.SellerID,
			&match.ProductID,
			&match.SizeID,
			&match.Price,
			&match.Quantity,
			&match.Status,
			&match.CompletedAt,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan match: %w", err)
		}
		matches = append(matches, match)
	}

	return matches, total, nil
}

// BeginTx starts a new transaction
func (r *BiddingRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
