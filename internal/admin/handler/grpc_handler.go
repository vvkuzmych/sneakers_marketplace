package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/middleware"
	adminpb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/admin"
)

// AdminHandler implements the gRPC AdminService interface
type AdminHandler struct {
	adminpb.UnimplementedAdminServiceServer
	service *service.AdminService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// ==================== User Management ====================

// ListUsers retrieves a paginated list of users
func (h *AdminHandler) ListUsers(ctx context.Context, req *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	params := model.ListUsersParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Role:     req.Role,
		Search:   req.Search,
	}

	users, total, err := h.service.ListUsers(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUsers := make([]*adminpb.User, len(users))
	for i, u := range users {
		pbUsers[i] = userToProto(&u)
	}

	return &adminpb.ListUsersResponse{
		Users:    pbUsers,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetUser retrieves a user with statistics
func (h *AdminHandler) GetUser(ctx context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, stats, err := h.service.GetUserWithStats(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &adminpb.GetUserResponse{
		User:       userToProto(user),
		Statistics: userStatsToProto(stats),
	}, nil
}

// BanUser bans a user
func (h *AdminHandler) BanUser(ctx context.Context, req *adminpb.BanUserRequest) (*adminpb.BanUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reason is required")
	}

	// Get admin ID from context
	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	// Get IP address from metadata
	ipAddress := getIPFromContext(ctx)

	err = h.service.BanUser(ctx, req.UserId, adminID, req.Reason, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.BanUserResponse{
		Success: true,
		Message: "User banned successfully",
	}, nil
}

// UnbanUser unbans a user
func (h *AdminHandler) UnbanUser(ctx context.Context, req *adminpb.UnbanUserRequest) (*adminpb.UnbanUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)

	err = h.service.UnbanUser(ctx, req.UserId, adminID, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.UnbanUserResponse{
		Success: true,
		Message: "User unbanned successfully",
	}, nil
}

// DeleteUser deletes a user
func (h *AdminHandler) DeleteUser(ctx context.Context, req *adminpb.DeleteUserRequest) (*adminpb.DeleteUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)
	reason := req.Reason
	if reason == "" {
		reason = "Deleted by admin"
	}

	err = h.service.DeleteUser(ctx, req.UserId, adminID, reason, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// UpdateUserRole updates a user's role
func (h *AdminHandler) UpdateUserRole(ctx context.Context, req *adminpb.UpdateUserRoleRequest) (*adminpb.UpdateUserRoleResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.NewRole == "" {
		return nil, status.Error(codes.InvalidArgument, "new_role is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)

	err = h.service.UpdateUserRole(ctx, req.UserId, adminID, req.NewRole, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.UpdateUserRoleResponse{
		Success: true,
		Message: "User role updated successfully",
	}, nil
}

// ==================== Order Management ====================

// ListAllOrders retrieves all orders
func (h *AdminHandler) ListAllOrders(ctx context.Context, req *adminpb.ListAllOrdersRequest) (*adminpb.ListAllOrdersResponse, error) {
	params := model.ListOrdersParams{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.DateFrom != nil {
		t := req.DateFrom.AsTime()
		params.DateFrom = &t
	}
	if req.DateTo != nil {
		t := req.DateTo.AsTime()
		params.DateTo = &t
	}

	orders, total, err := h.service.ListAllOrders(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbOrders := make([]*adminpb.OrderSummary, len(orders))
	for i, o := range orders {
		pbOrders[i] = orderSummaryToProto(&o)
	}

	return &adminpb.ListAllOrdersResponse{
		Orders:   pbOrders,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetOrderDetails retrieves order details
func (h *AdminHandler) GetOrderDetails(ctx context.Context, req *adminpb.GetOrderDetailsRequest) (*adminpb.GetOrderDetailsResponse, error) {
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, history, err := h.service.GetOrderDetails(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	pbHistory := make([]*adminpb.OrderStatusChange, len(history))
	for i, h := range history {
		pbHistory[i] = &adminpb.OrderStatusChange{
			Status:    h.Status,
			ChangedBy: h.ChangedBy,
			Notes:     h.Notes,
			ChangedAt: timestamppb.New(h.ChangedAt),
		}
	}

	return &adminpb.GetOrderDetailsResponse{
		Order:         orderSummaryToProto(order),
		StatusHistory: pbHistory,
	}, nil
}

// CancelOrder cancels an order
func (h *AdminHandler) CancelOrder(ctx context.Context, req *adminpb.CancelOrderRequest) (*adminpb.CancelOrderResponse, error) {
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)
	reason := req.Reason
	if reason == "" {
		reason = "Canceled by admin"
	}

	err = h.service.CancelOrder(ctx, req.OrderId, adminID, reason, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.CancelOrderResponse{
		Success: true,
		Message: "Order canceled successfully",
	}, nil
}

// ==================== Product Management ====================

// ListAllProducts retrieves all products
func (h *AdminHandler) ListAllProducts(ctx context.Context, req *adminpb.ListAllProductsRequest) (*adminpb.ListAllProductsResponse, error) {
	params := model.ListProductsParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Search:   req.Search,
	}

	products, total, err := h.service.ListAllProducts(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbProducts := make([]*adminpb.ProductSummary, len(products))
	for i, p := range products {
		pbProducts[i] = productSummaryToProto(&p)
	}

	return &adminpb.ListAllProductsResponse{
		Products: pbProducts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// FeatureProduct marks a product as featured
func (h *AdminHandler) FeatureProduct(ctx context.Context, req *adminpb.FeatureProductRequest) (*adminpb.FeatureProductResponse, error) {
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)

	err = h.service.FeatureProduct(ctx, req.ProductId, adminID, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.FeatureProductResponse{
		Success: true,
		Message: "Product featured successfully",
	}, nil
}

// HideProduct hides a product
func (h *AdminHandler) HideProduct(ctx context.Context, req *adminpb.HideProductRequest) (*adminpb.HideProductResponse, error) {
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	adminID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "admin not authenticated")
	}

	ipAddress := getIPFromContext(ctx)
	reason := req.Reason
	if reason == "" {
		reason = "Hidden by admin"
	}

	err = h.service.HideProduct(ctx, req.ProductId, adminID, reason, ipAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.HideProductResponse{
		Success: true,
		Message: "Product hidden successfully",
	}, nil
}

// ==================== Analytics ====================

// GetPlatformStats retrieves platform statistics
func (h *AdminHandler) GetPlatformStats(ctx context.Context, req *adminpb.GetPlatformStatsRequest) (*adminpb.GetPlatformStatsResponse, error) {
	stats, err := h.service.GetPlatformStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.GetPlatformStatsResponse{
		TotalUsers:         stats.TotalUsers,
		ActiveUsersToday:   stats.ActiveUsersToday,
		TotalProducts:      stats.TotalProducts,
		ActiveProducts:     stats.ActiveProducts,
		TotalOrders:        stats.TotalOrders,
		OrdersToday:        stats.OrdersToday,
		TotalRevenue:       stats.TotalRevenue,
		RevenueToday:       stats.RevenueToday,
		TotalFeesCollected: stats.TotalFeesCollected,
		TotalMatches:       stats.TotalMatches,
		MatchesToday:       stats.MatchesToday,
	}, nil
}

// GetRevenueReport retrieves revenue report
func (h *AdminHandler) GetRevenueReport(ctx context.Context, req *adminpb.GetRevenueReportRequest) (*adminpb.GetRevenueReportResponse, error) {
	if req.DateFrom == nil || req.DateTo == nil {
		return nil, status.Error(codes.InvalidArgument, "date_from and date_to are required")
	}

	params := model.GetRevenueReportParams{
		DateFrom: req.DateFrom.AsTime(),
		DateTo:   req.DateTo.AsTime(),
		GroupBy:  req.GroupBy,
	}

	report, err := h.service.GetRevenueReport(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dataPoints := make([]*adminpb.RevenueDataPoint, len(report.DataPoints))
	for i, dp := range report.DataPoints {
		dataPoints[i] = &adminpb.RevenueDataPoint{
			Label:      dp.Label,
			Revenue:    dp.Revenue,
			Fees:       dp.Fees,
			OrderCount: dp.OrderCount,
		}
	}

	return &adminpb.GetRevenueReportResponse{
		DataPoints:   dataPoints,
		TotalRevenue: report.TotalRevenue,
		TotalFees:    report.TotalFees,
	}, nil
}

// GetUserActivityReport retrieves user activity report
func (h *AdminHandler) GetUserActivityReport(ctx context.Context, req *adminpb.GetUserActivityReportRequest) (*adminpb.GetUserActivityReportResponse, error) {
	if req.DateFrom == nil || req.DateTo == nil {
		return nil, status.Error(codes.InvalidArgument, "date_from and date_to are required")
	}

	params := model.GetUserActivityReportParams{
		DateFrom: req.DateFrom.AsTime(),
		DateTo:   req.DateTo.AsTime(),
	}

	report, err := h.service.GetUserActivityReport(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &adminpb.GetUserActivityReportResponse{
		NewUsers:            report.NewUsers,
		ActiveUsers:         report.ActiveUsers,
		TotalBidsPlaced:     report.TotalBidsPlaced,
		TotalAsksPlaced:     report.TotalAsksPlaced,
		TotalMatchesCreated: report.TotalMatchesCreated,
	}, nil
}

// ==================== System Health ====================

// GetSystemHealth retrieves system health status
func (h *AdminHandler) GetSystemHealth(ctx context.Context, req *adminpb.GetSystemHealthRequest) (*adminpb.GetSystemHealthResponse, error) {
	// For now, return a simple healthy status
	// In production, this would check all services
	return &adminpb.GetSystemHealthResponse{
		OverallHealthy: true,
		Services:       []*adminpb.ServiceHealth{},
		Database:       &adminpb.DatabaseHealth{Healthy: true},
		CheckedAt:      timestamppb.Now(),
	}, nil
}

// GetServiceMetrics retrieves service metrics
func (h *AdminHandler) GetServiceMetrics(ctx context.Context, req *adminpb.GetServiceMetricsRequest) (*adminpb.GetServiceMetricsResponse, error) {
	// For now, return empty metrics
	// In production, this would collect metrics from all services
	return &adminpb.GetServiceMetricsResponse{
		Metrics: []*adminpb.ServiceMetric{},
	}, nil
}

// ==================== Audit Logs ====================

// GetAuditLogs retrieves audit logs
func (h *AdminHandler) GetAuditLogs(ctx context.Context, req *adminpb.GetAuditLogsRequest) (*adminpb.GetAuditLogsResponse, error) {
	params := model.ListAuditLogsParams{
		Page:       req.Page,
		PageSize:   req.PageSize,
		ActionType: req.ActionType,
		AdminID:    req.AdminId,
	}

	if req.DateFrom != nil {
		t := req.DateFrom.AsTime()
		params.DateFrom = &t
	}
	if req.DateTo != nil {
		t := req.DateTo.AsTime()
		params.DateTo = &t
	}

	logs, total, err := h.service.GetAuditLogs(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbLogs := make([]*adminpb.AuditLog, len(logs))
	for i, log := range logs {
		pbLogs[i] = auditLogToProto(&log)
	}

	return &adminpb.GetAuditLogsResponse{
		Logs:     pbLogs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ==================== Helper Functions: Model -> Proto ====================

func userToProto(u *model.AdminUser) *adminpb.User {
	user := &adminpb.User{
		Id:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Phone:       u.Phone,
		Role:        u.Role,
		IsActive:    u.IsActive,
		IsBanned:    u.IsBanned,
		TotalOrders: u.TotalOrders,
		TotalSpent:  u.TotalSpent,
		CreatedAt:   timestamppb.New(u.CreatedAt),
	}

	if u.BanReason != nil {
		user.BanReason = *u.BanReason
	}
	if u.BannedAt != nil {
		user.BannedAt = timestamppb.New(*u.BannedAt)
	}
	if u.LastLogin != nil {
		user.LastLogin = timestamppb.New(*u.LastLogin)
	}

	return user
}

func userStatsToProto(s *model.UserStatistics) *adminpb.UserStatistics {
	return &adminpb.UserStatistics{
		TotalBids:    s.TotalBids,
		TotalAsks:    s.TotalAsks,
		TotalMatches: s.TotalMatches,
		TotalOrders:  s.TotalOrders,
		TotalSpent:   s.TotalSpent,
		TotalEarned:  s.TotalEarned,
	}
}

func orderSummaryToProto(o *model.OrderSummary) *adminpb.OrderSummary {
	return &adminpb.OrderSummary{
		Id:          o.ID,
		OrderNumber: o.OrderNumber,
		BuyerId:     o.BuyerID,
		SellerId:    o.SellerID,
		BuyerEmail:  o.BuyerEmail,
		SellerEmail: o.SellerEmail,
		ProductId:   o.ProductID,
		ProductName: o.ProductName,
		Subtotal:    o.Subtotal,
		BuyerFee:    o.BuyerFee,
		SellerFee:   o.SellerFee,
		Total:       o.Total,
		Status:      o.Status,
		CreatedAt:   timestamppb.New(o.CreatedAt),
	}
}

func productSummaryToProto(p *model.ProductSummary) *adminpb.ProductSummary {
	return &adminpb.ProductSummary{
		Id:           p.ID,
		Sku:          p.SKU,
		Name:         p.Name,
		Brand:        p.Brand,
		Model:        p.Model,
		RetailPrice:  p.RetailPrice,
		IsActive:     p.IsActive,
		IsFeatured:   p.IsFeatured,
		TotalBids:    p.TotalBids,
		TotalAsks:    p.TotalAsks,
		TotalMatches: p.TotalMatches,
		HighestBid:   p.HighestBid,
		LowestAsk:    p.LowestAsk,
		CreatedAt:    timestamppb.New(p.CreatedAt),
	}
}

func auditLogToProto(log *model.AuditLog) *adminpb.AuditLog {
	detailsJSON, _ := log.DetailsJSON() // Ignore error, return empty string on failure
	return &adminpb.AuditLog{
		Id:         log.ID,
		AdminId:    log.AdminID,
		AdminEmail: log.AdminEmail,
		ActionType: log.ActionType,
		EntityType: log.EntityType,
		EntityId:   log.EntityID,
		Details:    detailsJSON,
		IpAddress:  log.IPAddress,
		CreatedAt:  timestamppb.New(log.CreatedAt),
	}
}

// getIPFromContext extracts IP address from gRPC metadata
func getIPFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "unknown"
	}

	// Try to get X-Forwarded-For first
	if xForwardedFor := md.Get("x-forwarded-for"); len(xForwardedFor) > 0 {
		return xForwardedFor[0]
	}

	// Try to get X-Real-IP
	if xRealIP := md.Get("x-real-ip"); len(xRealIP) > 0 {
		return xRealIP[0]
	}

	// Fallback
	return "127.0.0.1"
}
