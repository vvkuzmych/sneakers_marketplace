package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vvkuzmych/sneakers_marketplace/internal/order/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/order/repository"
)

type OrderService struct {
	repo *repository.OrderRepository

	// Fee percentages (can be loaded from config)
	defaultBuyerFeePercentage  float64
	defaultSellerFeePercentage float64
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo:                       repo,
		defaultBuyerFeePercentage:  0.03, // 3% buyer processing fee
		defaultSellerFeePercentage: 0.09, // 9% seller commission
	}
}

// CreateOrderFromMatch creates an order from a match
// This would typically be called by the Bidding Service after a match is created
func (s *OrderService) CreateOrderFromMatch(
	ctx context.Context,
	matchID int64,
	buyerID, sellerID int64,
	productID, sizeID int64,
	price float64,
	quantity int32,
	shippingAddressID *int64,
	buyerNotes string,
	buyerFeePercentage, sellerFeePercentage *float64,
) (*model.Order, error) {
	// Check if order already exists for this match
	existingOrder, err := s.repo.GetOrderByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing order: %w", err)
	}
	if existingOrder != nil {
		return existingOrder, nil // Order already exists
	}

	// Use default fee percentages if not provided
	buyerFee := s.defaultBuyerFeePercentage
	if buyerFeePercentage != nil {
		buyerFee = *buyerFeePercentage
	}

	sellerFee := s.defaultSellerFeePercentage
	if sellerFeePercentage != nil {
		sellerFee = *sellerFeePercentage
	}

	// Calculate fees
	buyerFeeAmount := price * buyerFee
	sellerFeeAmount := price * sellerFee
	platformFeeAmount := sellerFeeAmount // Platform earns from seller commission
	totalAmount := price + buyerFeeAmount
	sellerPayout := price - sellerFeeAmount

	order := &model.Order{
		MatchID:  matchID,
		BuyerID:  buyerID,
		SellerID: sellerID,

		ProductID: productID,
		SizeID:    sizeID,

		Price:        price,
		Quantity:     quantity,
		BuyerFee:     buyerFeeAmount,
		SellerFee:    sellerFeeAmount,
		PlatformFee:  platformFeeAmount,
		TotalAmount:  totalAmount,
		SellerPayout: sellerPayout,

		Status: model.StatusPendingPayment,
	}

	// Set shipping address if provided
	if shippingAddressID != nil {
		order.ShippingAddressID = sql.NullInt64{
			Int64: *shippingAddressID,
			Valid: true,
		}
	}

	// Set buyer notes if provided
	if buyerNotes != "" {
		order.BuyerNotes = sql.NullString{
			String: buyerNotes,
			Valid:  true,
		}
	}

	// Create order
	createdOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return createdOrder, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*model.Order, error) {
	return s.repo.GetOrderByID(ctx, orderID)
}

// GetOrderByNumber retrieves an order by order number
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	return s.repo.GetOrderByOrderNumber(ctx, orderNumber)
}

// ListOrders lists orders with optional filters
func (s *OrderService) ListOrders(ctx context.Context, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.ListOrders(ctx, status, page, pageSize)
}

// GetBuyerOrders retrieves orders for a buyer
func (s *OrderService) GetBuyerOrders(ctx context.Context, buyerID int64, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.GetBuyerOrders(ctx, buyerID, status, page, pageSize)
}

// GetSellerOrders retrieves orders for a seller
func (s *OrderService) GetSellerOrders(ctx context.Context, sellerID int64, status string, page, pageSize int32) ([]*model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.GetSellerOrders(ctx, sellerID, status, page, pageSize)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int64, newStatus, updatedBy string) (*model.Order, error) {
	// Validate new status
	if !model.IsValidStatus(newStatus) {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}

	// Get current order
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Check if transition is allowed
	if !order.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("cannot transition from %s to %s", order.Status, newStatus)
	}

	// Update status
	err = s.repo.UpdateOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		return nil, err
	}

	// Return updated order
	return s.repo.GetOrderByID(ctx, orderID)
}

// MarkAsPaid marks an order as paid (called by Payment Service)
func (s *OrderService) MarkAsPaid(ctx context.Context, orderID int64) (*model.Order, error) {
	return s.UpdateOrderStatus(ctx, orderID, model.StatusPaid, "payment_service")
}

// MarkAsProcessing marks an order as processing (seller acknowledges)
func (s *OrderService) MarkAsProcessing(ctx context.Context, orderID, sellerID int64) (*model.Order, error) {
	// Verify seller
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: not the seller of this order")
	}

	return s.UpdateOrderStatus(ctx, orderID, model.StatusProcessing, fmt.Sprintf("seller_%d", sellerID))
}

// AddTrackingNumber adds tracking information and marks as shipped
func (s *OrderService) AddTrackingNumber(ctx context.Context, orderID, sellerID int64, trackingNumber, carrier string) (*model.Order, error) {
	// Verify seller
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: not the seller of this order")
	}

	// Validate order status
	if order.Status != model.StatusProcessing && order.Status != model.StatusPaid {
		return nil, fmt.Errorf("cannot add tracking: order must be in processing or paid status")
	}

	// Validate inputs
	if trackingNumber == "" {
		return nil, fmt.Errorf("tracking number is required")
	}
	if carrier == "" {
		return nil, fmt.Errorf("carrier is required")
	}

	// Add tracking number (also updates status to shipped)
	err = s.repo.AddTrackingNumber(ctx, orderID, trackingNumber, carrier)
	if err != nil {
		return nil, err
	}

	return s.repo.GetOrderByID(ctx, orderID)
}

// MarkAsDelivered marks an order as delivered
func (s *OrderService) MarkAsDelivered(ctx context.Context, orderID int64) (*model.Order, error) {
	return s.UpdateOrderStatus(ctx, orderID, model.StatusDelivered, "system")
}

// MarkAsCompleted marks an order as completed
// This can be done by buyer confirmation or automatically after N days
func (s *OrderService) MarkAsCompleted(ctx context.Context, orderID int64, completedBy string) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != model.StatusDelivered {
		return nil, fmt.Errorf("order must be delivered before completion")
	}

	return s.UpdateOrderStatus(ctx, orderID, model.StatusCompleted, completedBy)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID int64, reason string) (*model.Order, error) {
	// Get order
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Check if user is authorized (buyer or seller)
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, fmt.Errorf("unauthorized: not a party to this order")
	}

	// Check if order can be cancelled
	if !order.CanBeCancelled() {
		return nil, fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	// Validate reason
	if reason == "" {
		return nil, fmt.Errorf("cancellation reason is required")
	}

	// Cancel order
	err = s.repo.CancelOrder(ctx, orderID, reason)
	if err != nil {
		return nil, err
	}

	// TODO: If order was paid, initiate refund via Payment Service

	return s.repo.GetOrderByID(ctx, orderID)
}

// GetOrderStatusHistory retrieves the status history for an order
func (s *OrderService) GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*model.OrderStatusHistory, error) {
	// Verify order exists
	_, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetOrderStatusHistory(ctx, orderID)
}

// GetShippingStatus gets shipping information for an order
func (s *OrderService) GetShippingStatus(ctx context.Context, orderID int64) (*model.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !order.TrackingNumber.Valid {
		return nil, fmt.Errorf("no tracking information available")
	}

	return order, nil
}
