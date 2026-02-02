package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/repository"
	feeService "github.com/vvkuzmych/sneakers_marketplace/internal/fees/service"
	notificationPb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/notification"
)

// BiddingService handles business logic for bidding
type BiddingService struct {
	repo               *repository.BiddingRepository
	notificationClient notificationPb.NotificationServiceClient
	feeService         *feeService.FeeService
}

// NewBiddingService creates a new bidding service
func NewBiddingService(repo *repository.BiddingRepository, notificationClient notificationPb.NotificationServiceClient, feeService *feeService.FeeService) *BiddingService {
	return &BiddingService{
		repo:               repo,
		notificationClient: notificationClient,
		feeService:         feeService,
	}
}

// PlaceBid places a new bid and attempts to match it
func (s *BiddingService) PlaceBid(ctx context.Context, userID, productID, sizeID int64, price float64, quantity int, expiresInHours int) (*model.Bid, *model.Match, error) {
	bid := &model.Bid{
		UserID:    userID,
		ProductID: productID,
		SizeID:    sizeID,
		Price:     price,
		Quantity:  quantity,
		Status:    model.StatusActive,
	}

	// Set expiration
	if expiresInHours > 0 {
		expiresAt := time.Now().Add(time.Duration(expiresInHours) * time.Hour)
		bid.ExpiresAt = &expiresAt
	}

	// Create bid
	if err := s.repo.PlaceBid(ctx, bid); err != nil {
		return nil, nil, fmt.Errorf("failed to place bid: %w", err)
	}

	// Try to match immediately
	match, err := s.tryMatchBid(ctx, bid)
	if err != nil {
		return bid, nil, fmt.Errorf("bid placed but matching failed: %w", err)
	}

	return bid, match, nil
}

// PlaceAsk places a new ask and attempts to match it
func (s *BiddingService) PlaceAsk(ctx context.Context, userID, productID, sizeID int64, price float64, quantity int, expiresInHours int) (*model.Ask, *model.Match, error) {
	ask := &model.Ask{
		UserID:    userID,
		ProductID: productID,
		SizeID:    sizeID,
		Price:     price,
		Quantity:  quantity,
		Status:    model.StatusActive,
	}

	// Set expiration
	if expiresInHours > 0 {
		expiresAt := time.Now().Add(time.Duration(expiresInHours) * time.Hour)
		ask.ExpiresAt = &expiresAt
	}

	// Create ask
	if err := s.repo.PlaceAsk(ctx, ask); err != nil {
		return nil, nil, fmt.Errorf("failed to place ask: %w", err)
	}

	// Try to match immediately
	match, err := s.tryMatchAsk(ctx, ask)
	if err != nil {
		return ask, nil, fmt.Errorf("ask placed but matching failed: %w", err)
	}

	return ask, match, nil
}

// tryMatchBid attempts to match a bid with existing asks
func (s *BiddingService) tryMatchBid(ctx context.Context, bid *model.Bid) (*model.Match, error) {
	// Find lowest ask that can be matched
	lowestAsk, err := s.repo.GetLowestAsk(ctx, bid.ProductID, bid.SizeID)
	if err != nil {
		return nil, err
	}

	// No asks available
	if lowestAsk == nil {
		return nil, nil
	}

	// Check if match is possible
	if !model.CanMatch(bid, lowestAsk) {
		return nil, nil
	}

	// Create match in transaction
	match, err := s.createMatch(ctx, bid, lowestAsk)
	if err != nil {
		return nil, err
	}

	return match, nil
}

// tryMatchAsk attempts to match an ask with existing bids
func (s *BiddingService) tryMatchAsk(ctx context.Context, ask *model.Ask) (*model.Match, error) {
	// Find highest bid that can be matched
	highestBid, err := s.repo.GetHighestBid(ctx, ask.ProductID, ask.SizeID)
	if err != nil {
		return nil, err
	}

	// No bids available
	if highestBid == nil {
		return nil, nil
	}

	// Check if match is possible
	if !model.CanMatch(highestBid, ask) {
		return nil, nil
	}

	// Create match in transaction
	match, err := s.createMatch(ctx, highestBid, ask)
	if err != nil {
		return nil, err
	}

	return match, nil
}

// createMatch creates a match between bid and ask
func (s *BiddingService) createMatch(ctx context.Context, bid *model.Bid, ask *model.Ask) (*model.Match, error) {
	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Match price is the ask price (seller's price)
	matchPrice := ask.Price

	// Create match record
	match := &model.Match{
		BidID:     bid.ID,
		AskID:     ask.ID,
		BuyerID:   bid.UserID,
		SellerID:  ask.UserID,
		ProductID: bid.ProductID,
		SizeID:    bid.SizeID,
		Price:     matchPrice,
		Quantity:  bid.Quantity, // Assuming equal quantities
		Status:    model.StatusPending,
	}

	if err := s.repo.CreateMatch(ctx, tx, match); err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	// Update bid status to matched
	if err := s.repo.UpdateBidStatus(ctx, tx, bid.ID, model.StatusMatched); err != nil {
		return nil, fmt.Errorf("failed to update bid status: %w", err)
	}

	// Update ask status to matched
	if err := s.repo.UpdateAskStatus(ctx, tx, ask.ID, model.StatusMatched); err != nil {
		return nil, fmt.Errorf("failed to update ask status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Calculate and record fees (after transaction committed)
	if s.feeService != nil {
		// TODO: Get vertical from product (for now assume sneakers)
		vertical := "sneakers"

		// Calculate fees based on seller's subscription tier
		feeBreakdown, err := s.feeService.CalculateFees(ctx, vertical, matchPrice, match.SellerID)
		if err != nil {
			log.Printf("Failed to calculate fees for match %d: %v", match.ID, err)
		} else {
			// Record transaction fee
			if err := s.feeService.RecordTransactionFee(ctx, match.ID, nil, feeBreakdown, vertical); err != nil {
				log.Printf("Failed to record transaction fee for match %d: %v", match.ID, err)
			}
		}
	}

	// Send notification asynchronously (don't block the response)
	if s.notificationClient != nil {
		go func() {
			notifyCtx := context.Background()
			_, err := s.notificationClient.NotifyMatchCreated(notifyCtx, &notificationPb.NotifyMatchCreatedRequest{
				MatchId:     match.ID,
				BuyerId:     match.BuyerID,
				SellerId:    match.SellerID,
				ProductId:   match.ProductID,
				ProductName: "Product", // TODO: Get product name from Product Service
				Size:        "10",      // TODO: Get size from Size table
				Price:       match.Price,
			})
			if err != nil {
				log.Printf("Failed to send match notification: %v", err)
			}
		}()
	}

	return match, nil
}

// GetBid retrieves a bid by ID
func (s *BiddingService) GetBid(ctx context.Context, bidID int64) (*model.Bid, error) {
	return s.repo.GetBidByID(ctx, bidID)
}

// GetUserBids retrieves bids for a user
func (s *BiddingService) GetUserBids(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Bid, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetUserBids(ctx, userID, status, page, pageSize)
}

// GetProductBids retrieves bids for a product
func (s *BiddingService) GetProductBids(ctx context.Context, productID, sizeID int64, status string, page, pageSize int) ([]*model.Bid, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetProductBids(ctx, productID, sizeID, status, page, pageSize)
}

// CancelBid cancels a bid
func (s *BiddingService) CancelBid(ctx context.Context, bidID, userID int64) error {
	// Get bid
	bid, err := s.repo.GetBidByID(ctx, bidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", err)
	}

	// Check ownership
	if bid.UserID != userID {
		return fmt.Errorf("unauthorized: not bid owner")
	}

	// Check if can be canceled
	if bid.Status != model.StatusActive {
		return fmt.Errorf("bid cannot be canceled: status is %s", bid.Status)
	}

	// Update status
	return s.repo.UpdateBidStatus(ctx, nil, bidID, model.StatusCancelled)
}

// GetAsk retrieves an ask by ID
func (s *BiddingService) GetAsk(ctx context.Context, askID int64) (*model.Ask, error) {
	return s.repo.GetAskByID(ctx, askID)
}

// GetUserAsks retrieves asks for a user
func (s *BiddingService) GetUserAsks(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.Ask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetUserAsks(ctx, userID, status, page, pageSize)
}

// GetProductAsks retrieves asks for a product
func (s *BiddingService) GetProductAsks(ctx context.Context, productID, sizeID int64, status string, page, pageSize int) ([]*model.Ask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetProductAsks(ctx, productID, sizeID, status, page, pageSize)
}

// CancelAsk cancels an ask
func (s *BiddingService) CancelAsk(ctx context.Context, askID, userID int64) error {
	// Get ask
	ask, err := s.repo.GetAskByID(ctx, askID)
	if err != nil {
		return fmt.Errorf("ask not found: %w", err)
	}

	// Check ownership
	if ask.UserID != userID {
		return fmt.Errorf("unauthorized: not ask owner")
	}

	// Check if can be canceled
	if ask.Status != model.StatusActive {
		return fmt.Errorf("ask cannot be canceled: status is %s", ask.Status)
	}

	// Update status
	return s.repo.UpdateAskStatus(ctx, nil, askID, model.StatusCancelled)
}

// GetHighestBid retrieves the highest bid for a product/size
func (s *BiddingService) GetHighestBid(ctx context.Context, productID, sizeID int64) (*model.Bid, error) {
	return s.repo.GetHighestBid(ctx, productID, sizeID)
}

// GetLowestAsk retrieves the lowest ask for a product/size
func (s *BiddingService) GetLowestAsk(ctx context.Context, productID, sizeID int64) (*model.Ask, error) {
	return s.repo.GetLowestAsk(ctx, productID, sizeID)
}

// GetMarketPrice retrieves market data for a product/size
func (s *BiddingService) GetMarketPrice(ctx context.Context, productID, sizeID int64) (*model.MarketPrice, error) {
	marketPrice := &model.MarketPrice{
		ProductID: productID,
		SizeID:    sizeID,
	}

	// Get highest bid
	highestBid, err := s.repo.GetHighestBid(ctx, productID, sizeID)
	if err == nil && highestBid != nil {
		marketPrice.HighestBid = highestBid.Price
	}

	// Get lowest ask
	lowestAsk, err := s.repo.GetLowestAsk(ctx, productID, sizeID)
	if err == nil && lowestAsk != nil {
		marketPrice.LowestAsk = lowestAsk.Price
	}

	// Get total counts
	bids, totalBids, err := s.repo.GetProductBids(ctx, productID, sizeID, model.StatusActive, 1, 1)
	if err == nil {
		marketPrice.TotalBids = totalBids
	}
	_ = bids // unused

	asks, totalAsks, err := s.repo.GetProductAsks(ctx, productID, sizeID, model.StatusActive, 1, 1)
	if err == nil {
		marketPrice.TotalAsks = totalAsks
	}
	_ = asks // unused

	// TODO: Get last sale price from matches table
	marketPrice.LastSale = 0

	return marketPrice, nil
}

// GetMatch retrieves a match by ID
func (s *BiddingService) GetMatch(ctx context.Context, matchID int64) (*model.Match, error) {
	return s.repo.GetMatchByID(ctx, matchID)
}

// GetUserMatches retrieves matches for a user
func (s *BiddingService) GetUserMatches(ctx context.Context, userID int64, asBuyer, asSeller bool, page, pageSize int) ([]*model.Match, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.GetUserMatches(ctx, userID, asBuyer, asSeller, page, pageSize)
}
