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
	repo *repository.FeeRepository
	log  *logger.Logger
}

func NewFeeService(repo *repository.FeeRepository, log *logger.Logger) *FeeService {
	return &FeeService{
		repo: repo,
		log:  log,
	}
}

// CalculateFees calculates all fees for a transaction
// Returns detailed fee breakdown for buyer and seller
func (s *FeeService) CalculateFees(ctx context.Context, vertical string, salePrice float64, includeAuth bool) (*model.FeeBreakdown, error) {
	// Validate inputs
	if salePrice <= 0 {
		return nil, fmt.Errorf("sale price must be positive, got: %.2f", salePrice)
	}

	if vertical == "" {
		vertical = "sneakers" // default
	}

	// Get fee config for vertical
	config, err := s.repo.GetFeeConfig(ctx, vertical)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee config for vertical %s: %w", vertical, err)
	}

	// Create breakdown
	breakdown := model.NewFeeBreakdown(salePrice)

	// NEW MODEL: Buyer pays transaction fee, Seller receives full price
	transactionFee := config.CalculateTransactionFee(salePrice)

	// Buyer fees (transaction fee instead of seller)
	breakdown.BuyerProcessingFee = s.roundToTwoDecimals(transactionFee)
	breakdown.BuyerShippingFee = 0.0 // No shipping fees

	// Seller fees (zero - seller receives full price)
	breakdown.SellerTransactionFee = 0.0
	breakdown.SellerAuthFee = 0.0
	breakdown.SellerShippingCost = 0.0

	// Calculate totals
	breakdown.CalculateTotals()

	s.log.Infof("Fee calculation for %s @ $%.2f: Platform Revenue = $%.2f (Buyer: $%.2f, Seller pays: $%.2f)",
		vertical,
		salePrice,
		breakdown.PlatformRevenue,
		breakdown.BuyerTotal,
		breakdown.SellerTransactionFee+breakdown.SellerAuthFee+breakdown.SellerShippingCost,
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
