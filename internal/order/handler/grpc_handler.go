package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vvkuzmych/sneakers_marketplace/internal/order/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/order/service"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/order"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// Helper function to convert model.Order to pb.Order
func orderToProto(o *model.Order) *pb.Order {
	if o == nil {
		return nil
	}

	order := &pb.Order{
		Id:           o.ID,
		OrderNumber:  o.OrderNumber,
		MatchId:      o.MatchID,
		BuyerId:      o.BuyerID,
		SellerId:     o.SellerID,
		ProductId:    o.ProductID,
		SizeId:       o.SizeID,
		Price:        o.Price,
		Quantity:     o.Quantity,
		BuyerFee:     o.BuyerFee,
		SellerFee:    o.SellerFee,
		PlatformFee:  o.PlatformFee,
		TotalAmount:  o.TotalAmount,
		SellerPayout: o.SellerPayout,
		Status:       o.Status,
		CreatedAt:    timestamppb.New(o.CreatedAt),
		UpdatedAt:    timestamppb.New(o.UpdatedAt),
	}

	// Optional fields
	if o.ShippingAddressID.Valid {
		order.ShippingAddressId = o.ShippingAddressID.Int64
	}
	if o.TrackingNumber.Valid {
		order.TrackingNumber = o.TrackingNumber.String
	}
	if o.Carrier.Valid {
		order.Carrier = o.Carrier.String
	}
	if o.PaymentAt.Valid {
		order.PaymentAt = timestamppb.New(o.PaymentAt.Time)
	}
	if o.ShippedAt.Valid {
		order.ShippedAt = timestamppb.New(o.ShippedAt.Time)
	}
	if o.DeliveredAt.Valid {
		order.DeliveredAt = timestamppb.New(o.DeliveredAt.Time)
	}
	if o.CompletedAt.Valid {
		order.CompletedAt = timestamppb.New(o.CompletedAt.Time)
	}
	if o.CancelledAt.Valid {
		order.CancelledAt = timestamppb.New(o.CancelledAt.Time)
	}
	if o.BuyerNotes.Valid {
		order.BuyerNotes = o.BuyerNotes.String
	}
	if o.SellerNotes.Valid {
		order.SellerNotes = o.SellerNotes.String
	}
	if o.AdminNotes.Valid {
		order.AdminNotes = o.AdminNotes.String
	}
	if o.CancellationReason.Valid {
		order.CancellationReason = o.CancellationReason.String
	}

	return order
}

// Helper function to convert model.OrderStatusHistory to pb.OrderStatusHistory
func statusHistoryToProto(h *model.OrderStatusHistory) *pb.OrderStatusHistory {
	if h == nil {
		return nil
	}

	history := &pb.OrderStatusHistory{
		Id:        h.ID,
		OrderId:   h.OrderID,
		ToStatus:  h.ToStatus,
		CreatedBy: h.CreatedBy,
		CreatedAt: timestamppb.New(h.CreatedAt),
	}

	if h.FromStatus.Valid {
		history.FromStatus = h.FromStatus.String
	}
	if h.Note.Valid {
		history.Note = h.Note.String
	}

	return history
}

// CreateOrder creates a new order from a match
func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if req.MatchId == 0 {
		return &pb.CreateOrderResponse{
			Error: "match_id is required",
		}, nil
	}

	// TODO: Fetch match details from Bidding Service to get buyer, seller, product, etc.
	// For now, we'll require these to be passed somehow or fetched
	// This is a simplified version

	// Note: In a real implementation, we'd fetch match details here
	// For now, return error asking for implementation
	return &pb.CreateOrderResponse{
		Error: "CreateOrder requires integration with Bidding Service to fetch match details",
	}, nil
}

// GetOrder retrieves an order by ID
func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if req.OrderId == 0 {
		return &pb.GetOrderResponse{
			Error: "order_id is required",
		}, nil
	}

	order, err := h.service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return &pb.GetOrderResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetOrderResponse{
		Order: orderToProto(order),
	}, nil
}

// ListOrders lists all orders with optional filters
func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	page := req.Page
	if page == 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	orders, total, err := h.service.ListOrders(ctx, req.Status, page, pageSize)
	if err != nil {
		return &pb.ListOrdersResponse{
			Error: err.Error(),
		}, nil
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = orderToProto(order)
	}

	return &pb.ListOrdersResponse{
		Orders:   pbOrders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	if req.OrderId == 0 {
		return &pb.UpdateOrderStatusResponse{
			Error: "order_id is required",
		}, nil
	}
	if req.NewStatus == "" {
		return &pb.UpdateOrderStatusResponse{
			Error: "new_status is required",
		}, nil
	}

	updatedBy := req.UpdatedBy
	if updatedBy == "" {
		updatedBy = "system"
	}

	order, err := h.service.UpdateOrderStatus(ctx, req.OrderId, req.NewStatus, updatedBy)
	if err != nil {
		return &pb.UpdateOrderStatusResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateOrderStatusResponse{
		Order: orderToProto(order),
	}, nil
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	if req.OrderId == 0 {
		return &pb.CancelOrderResponse{
			Error: "order_id is required",
		}, nil
	}
	if req.CancelledByUserId == 0 {
		return &pb.CancelOrderResponse{
			Error: "canceled_by_user_id is required",
		}, nil
	}
	if req.Reason == "" {
		return &pb.CancelOrderResponse{
			Error: "cancellation reason is required",
		}, nil
	}

	order, err := h.service.CancelOrder(ctx, req.OrderId, req.CancelledByUserId, req.Reason)
	if err != nil {
		return &pb.CancelOrderResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.CancelOrderResponse{
		Order: orderToProto(order),
	}, nil
}

// AddTrackingNumber adds tracking information to an order
func (h *OrderHandler) AddTrackingNumber(ctx context.Context, req *pb.AddTrackingNumberRequest) (*pb.AddTrackingNumberResponse, error) {
	if req.OrderId == 0 {
		return &pb.AddTrackingNumberResponse{
			Error: "order_id is required",
		}, nil
	}
	if req.SellerId == 0 {
		return &pb.AddTrackingNumberResponse{
			Error: "seller_id is required",
		}, nil
	}
	if req.TrackingNumber == "" {
		return &pb.AddTrackingNumberResponse{
			Error: "tracking_number is required",
		}, nil
	}
	if req.Carrier == "" {
		return &pb.AddTrackingNumberResponse{
			Error: "carrier is required",
		}, nil
	}

	order, err := h.service.AddTrackingNumber(ctx, req.OrderId, req.SellerId, req.TrackingNumber, req.Carrier)
	if err != nil {
		return &pb.AddTrackingNumberResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.AddTrackingNumberResponse{
		Order: orderToProto(order),
	}, nil
}

// GetShippingStatus gets shipping status for an order
func (h *OrderHandler) GetShippingStatus(ctx context.Context, req *pb.GetShippingStatusRequest) (*pb.GetShippingStatusResponse, error) {
	if req.OrderId == 0 {
		return &pb.GetShippingStatusResponse{
			Error: "order_id is required",
		}, nil
	}

	order, err := h.service.GetShippingStatus(ctx, req.OrderId)
	if err != nil {
		return &pb.GetShippingStatusResponse{
			Error: err.Error(),
		}, nil
	}

	response := &pb.GetShippingStatusResponse{
		Status: order.Status,
	}

	if order.TrackingNumber.Valid {
		response.TrackingNumber = order.TrackingNumber.String
	}
	if order.Carrier.Valid {
		response.Carrier = order.Carrier.String
	}
	if order.ShippedAt.Valid {
		response.ShippedAt = timestamppb.New(order.ShippedAt.Time)
	}
	if order.DeliveredAt.Valid {
		response.DeliveredAt = timestamppb.New(order.DeliveredAt.Time)
	}

	return response, nil
}

// GetBuyerOrders retrieves orders for a buyer
func (h *OrderHandler) GetBuyerOrders(ctx context.Context, req *pb.GetBuyerOrdersRequest) (*pb.GetBuyerOrdersResponse, error) {
	if req.BuyerId == 0 {
		return &pb.GetBuyerOrdersResponse{
			Error: "buyer_id is required",
		}, nil
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	orders, total, err := h.service.GetBuyerOrders(ctx, req.BuyerId, req.Status, page, pageSize)
	if err != nil {
		return &pb.GetBuyerOrdersResponse{
			Error: err.Error(),
		}, nil
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = orderToProto(order)
	}

	return &pb.GetBuyerOrdersResponse{
		Orders: pbOrders,
		Total:  total,
	}, nil
}

// GetSellerOrders retrieves orders for a seller
func (h *OrderHandler) GetSellerOrders(ctx context.Context, req *pb.GetSellerOrdersRequest) (*pb.GetSellerOrdersResponse, error) {
	if req.SellerId == 0 {
		return &pb.GetSellerOrdersResponse{
			Error: "seller_id is required",
		}, nil
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	orders, total, err := h.service.GetSellerOrders(ctx, req.SellerId, req.Status, page, pageSize)
	if err != nil {
		return &pb.GetSellerOrdersResponse{
			Error: err.Error(),
		}, nil
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = orderToProto(order)
	}

	return &pb.GetSellerOrdersResponse{
		Orders: pbOrders,
		Total:  total,
	}, nil
}

// GetOrderStatusHistory retrieves the status history for an order
func (h *OrderHandler) GetOrderStatusHistory(ctx context.Context, req *pb.GetOrderStatusHistoryRequest) (*pb.GetOrderStatusHistoryResponse, error) {
	if req.OrderId == 0 {
		return &pb.GetOrderStatusHistoryResponse{
			Error: "order_id is required",
		}, nil
	}

	history, err := h.service.GetOrderStatusHistory(ctx, req.OrderId)
	if err != nil {
		return &pb.GetOrderStatusHistoryResponse{
			Error: err.Error(),
		}, nil
	}

	pbHistory := make([]*pb.OrderStatusHistory, len(history))
	for i, h := range history {
		pbHistory[i] = statusHistoryToProto(h)
	}

	return &pb.GetOrderStatusHistoryResponse{
		History: pbHistory,
	}, nil
}
