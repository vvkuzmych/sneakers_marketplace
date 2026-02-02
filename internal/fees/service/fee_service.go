package service

import (
	"context"
	"fmt"
	"math"

	"github.com/vvkuzmych/sneakers_marketplace/internal/fees/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/fees/repository"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

type FeeService struct {
	repo                 *repository.FeeRepository
	log                  *logger.Logger
	subscriptionProvider SubscriptionFeeProvider
}

func NewFeeService(repo *repository.FeeRepository, log *logger.Logger, subscriptionProvider SubscriptionFeeProvider) *FeeService {
	// If no subscription provider, use default
	if subscriptionProvider == nil {
		subscriptionProvider = NewDefaultFeeProvider()
	}

	return &FeeService{
		repo:                 repo,
		log:                  log,
		subscriptionProvider: subscriptionProvider,
	}
}

// CalculateFees calculates all fees for a transaction with subscription-based pricing
// Returns detailed fee breakdown for buyer and seller
// sellerUserID is used to determine seller's subscription tier and apply appropriate fees
// Use sellerUserID = -1 for preview/default fees (Free tier: 1%)
func (s *FeeService) CalculateFees(ctx context.Context, vertical string, salePrice float64, sellerUserID int64) (*model.FeeBreakdown, error) {
	// Validate inputs
	if salePrice <= 0 {
		return nil, fmt.Errorf("sale price must be positive, got: %.2f", salePrice)
	}

	if vertical == "" {
		vertical = "sneakers" // default
	}

	// Handle special case: -1 means use default fees (for preview)
	var sellerFeePercent, buyerFeePercent float64
	if sellerUserID == -1 {
		// Default Free tier fees
		sellerFeePercent = 1.0
		buyerFeePercent = 1.0
		s.log.Infof("Using default fees (preview mode): seller_fee=%.2f%%, buyer_fee=%.2f%%",
			sellerFeePercent, buyerFeePercent)
	} else if sellerUserID <= 0 {
		return nil, fmt.Errorf("sellerUserID must be positive or -1 for default fees, got: %d", sellerUserID)
	} else {
		// Get seller's subscription fee percentages
		var err error
		sellerFeePercent, buyerFeePercent, err = s.subscriptionProvider.GetUserFeePercentages(ctx, sellerUserID)
		if err != nil {
			s.log.Warnf("Failed to get subscription fees for seller %d, using defaults: %v", sellerUserID, err)
			// Fallback to default fees
			sellerFeePercent = 1.0
			buyerFeePercent = 1.0
		}

		s.log.Infof("Calculating fees for seller %d: seller_fee=%.2f%%, buyer_fee=%.2f%%",
			sellerUserID, sellerFeePercent, buyerFeePercent)
	}

	// Create breakdown
	breakdown := model.NewFeeBreakdown(salePrice)

	// SUBSCRIPTION-BASED MODEL:
	// - Seller pays platform fee based on their subscription tier (Free: 1%, Pro: 0.75%, Elite: 0.5%)
	// - Buyer pays fixed processing fee (1%)

	// Seller fees (based on subscription)
	sellerPlatformFee := s.roundToTwoDecimals((salePrice * sellerFeePercent) / 100)
	breakdown.SellerTransactionFee = sellerPlatformFee
	breakdown.SellerAuthFee = 0.0
	breakdown.SellerShippingCost = 0.0

	// Buyer fees (fixed processing)
	buyerProcessingFee := s.roundToTwoDecimals((salePrice * buyerFeePercent) / 100)
	breakdown.BuyerProcessingFee = buyerProcessingFee
	breakdown.BuyerShippingFee = 0.0

	// Calculate totals
	breakdown.CalculateTotals()

	s.log.Infof("Fee calculation for %s @ $%.2f: Seller pays $%.2f (%.2f%%), Buyer pays $%.2f (%.2f%%), Platform Revenue = $%.2f",
		vertical,
		salePrice,
		breakdown.SellerTransactionFee,
		sellerFeePercent,
		breakdown.BuyerProcessingFee,
		buyerFeePercent,
		breakdown.PlatformRevenue,
	)

	return breakdown, nil
}

// RecordTransactionFee saves fee record to database
// Should be called after a match is created
func (s *FeeService) RecordTransactionFee(ctx context.Context, matchID int64, orderID *int64, breakdown *model.FeeBreakdown, vertical string) error {
	// Get current config for snapshot
	config, err := s.repo.GetFeeConfig(ctx, vertical)
	if err != nil {
		return fmt.Errorf("failed to get fee config: %w", err)
	}

	// Convert breakdown to transaction fee model
	transFee := breakdown.ToTransactionFee(matchID, vertical, config)
	transFee.OrderID = orderID

	// Save to database
	if err := s.repo.CreateTransactionFee(ctx, transFee); err != nil {
		return fmt.Errorf("failed to record transaction fee: %w", err)
	}

	s.log.Infof("Transaction fee recorded: MatchID=%d, Revenue=$%.2f", matchID, breakdown.PlatformRevenue)

	return nil
}

// GetTransactionFee retrieves transaction fee by match ID
func (s *FeeService) GetTransactionFee(ctx context.Context, matchID int64) (*model.TransactionFee, error) {
	fee, err := s.repo.GetTransactionFeeByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction fee: %w", err)
	}
	return fee, nil
}

// GetFeeConfig retrieves fee configuration for a vertical
func (s *FeeService) GetFeeConfig(ctx context.Context, vertical string) (*model.FeeConfig, error) {
	return s.repo.GetFeeConfig(ctx, vertical)
}

// GetAllFeeConfigs retrieves all fee configurations
func (s *FeeService) GetAllFeeConfigs(ctx context.Context) ([]*model.FeeConfig, error) {
	return s.repo.GetAllFeeConfigs(ctx)
}

// roundToTwoDecimals rounds a float to 2 decimal places
func (s *FeeService) roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
