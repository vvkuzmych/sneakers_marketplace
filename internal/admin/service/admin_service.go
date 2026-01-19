package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/repository"
)

// AdminService handles business logic for admin operations
type AdminService struct {
	repo *repository.AdminRepository
}

// NewAdminService creates a new admin service
func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// ==================== User Management ====================

// ListUsers retrieves a list of users
func (s *AdminService) ListUsers(ctx context.Context, params model.ListUsersParams) ([]model.AdminUser, int32, error) {
	// Validate params
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	return s.repo.ListUsers(ctx, params)
}

// GetUserWithStats retrieves a user with statistics
func (s *AdminService) GetUserWithStats(ctx context.Context, userID int64) (*model.AdminUser, *model.UserStatistics, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	stats, err := s.repo.GetUserStatistics(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return user, stats, nil
}

// BanUser bans a user and creates audit log
func (s *AdminService) BanUser(ctx context.Context, userID, adminID int64, reason, ipAddress string) error {
	// Check if user exists and is not already banned
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsBanned {
		return fmt.Errorf("user is already banned")
	}

	// Don't allow banning admin users
	if user.Role == "admin" {
		return fmt.Errorf("cannot ban admin users")
	}

	// Ban user
	params := model.BanUserParams{
		UserID:   userID,
		AdminID:  adminID,
		Reason:   reason,
		BannedAt: time.Now(),
	}

	err = s.repo.BanUser(ctx, params)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionUserBanned,
		EntityType: model.EntityUser,
		EntityID:   userID,
		Details: map[string]interface{}{
			"reason":     reason,
			"user_email": user.Email,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// UnbanUser unbans a user and creates audit log
func (s *AdminService) UnbanUser(ctx context.Context, userID, adminID int64, ipAddress string) error {
	// Check if user exists and is banned
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsBanned {
		return fmt.Errorf("user is not banned")
	}

	// Unban user
	params := model.UnbanUserParams{
		UserID:  userID,
		AdminID: adminID,
	}

	err = s.repo.UnbanUser(ctx, params)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionUserUnbanned,
		EntityType: model.EntityUser,
		EntityID:   userID,
		Details: map[string]interface{}{
			"user_email": user.Email,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// DeleteUser deletes a user and creates audit log
func (s *AdminService) DeleteUser(ctx context.Context, userID, adminID int64, reason, ipAddress string) error {
	// Check if user exists
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user can be deleted
	if !user.CanBeDeleted() {
		return fmt.Errorf("cannot delete admin users with orders")
	}

	// Delete user
	params := model.DeleteUserParams{
		UserID:  userID,
		AdminID: adminID,
		Reason:  reason,
	}

	err = s.repo.DeleteUser(ctx, params)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionUserDeleted,
		EntityType: model.EntityUser,
		EntityID:   userID,
		Details: map[string]interface{}{
			"reason":     reason,
			"user_email": user.Email,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// UpdateUserRole updates a user's role and creates audit log
func (s *AdminService) UpdateUserRole(ctx context.Context, userID, adminID int64, newRole, ipAddress string) error {
	// Validate role
	if !model.IsValidRole(newRole) {
		return fmt.Errorf("invalid role: must be 'user' or 'admin'")
	}

	// Get user
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role == newRole {
		return fmt.Errorf("user already has role: %s", newRole)
	}

	// Update role
	params := model.UpdateUserRoleParams{
		UserID:  userID,
		AdminID: adminID,
		NewRole: newRole,
	}

	err = s.repo.UpdateUserRole(ctx, params)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionUserRoleUpdated,
		EntityType: model.EntityUser,
		EntityID:   userID,
		Details: map[string]interface{}{
			"old_role":   user.Role,
			"new_role":   newRole,
			"user_email": user.Email,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// ==================== Analytics ====================

// GetPlatformStats retrieves platform statistics
func (s *AdminService) GetPlatformStats(ctx context.Context) (*model.PlatformStats, error) {
	return s.repo.GetPlatformStats(ctx)
}

// GetRevenueReport retrieves revenue report
func (s *AdminService) GetRevenueReport(ctx context.Context, params model.GetRevenueReportParams) (*model.RevenueReport, error) {
	// Validate params
	if params.GroupBy == "" {
		params.GroupBy = "day"
	}

	if !isValidGroupBy(params.GroupBy) {
		return nil, fmt.Errorf("invalid groupBy: must be 'day', 'week', or 'month'")
	}

	return s.repo.GetRevenueReport(ctx, params)
}

// GetUserActivityReport retrieves user activity report
func (s *AdminService) GetUserActivityReport(ctx context.Context, params model.GetUserActivityReportParams) (*model.UserActivityReport, error) {
	return s.repo.GetUserActivityReport(ctx, params)
}

// ==================== Order Management ====================

// ListAllOrders retrieves all orders
func (s *AdminService) ListAllOrders(ctx context.Context, params model.ListOrdersParams) ([]model.OrderSummary, int32, error) {
	// Validate params
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	return s.repo.ListAllOrders(ctx, params)
}

// GetOrderDetails retrieves order details
func (s *AdminService) GetOrderDetails(ctx context.Context, orderID int64) (*model.OrderSummary, []model.OrderStatusChange, error) {
	return s.repo.GetOrderDetails(ctx, orderID)
}

// CancelOrder cancels an order and creates audit log
func (s *AdminService) CancelOrder(ctx context.Context, orderID, adminID int64, reason, ipAddress string) error {
	// Get order first
	order, _, err := s.repo.GetOrderDetails(ctx, orderID)
	if err != nil {
		return err
	}

	// Cancel order
	err = s.repo.CancelOrder(ctx, orderID, reason)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionOrderCancelled,
		EntityType: model.EntityOrder,
		EntityID:   orderID,
		Details: map[string]interface{}{
			"reason":       reason,
			"order_number": order.OrderNumber,
			"buyer_email":  order.BuyerEmail,
			"seller_email": order.SellerEmail,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// ==================== Product Management ====================

// ListAllProducts retrieves all products
func (s *AdminService) ListAllProducts(ctx context.Context, params model.ListProductsParams) ([]model.ProductSummary, int32, error) {
	// Validate params
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	return s.repo.ListAllProducts(ctx, params)
}

// FeatureProduct features a product and creates audit log
func (s *AdminService) FeatureProduct(ctx context.Context, productID, adminID int64, ipAddress string) error {
	err := s.repo.FeatureProduct(ctx, productID)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionProductFeatured,
		EntityType: model.EntityProduct,
		EntityID:   productID,
		Details:    map[string]interface{}{},
		IPAddress:  ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// HideProduct hides a product and creates audit log
func (s *AdminService) HideProduct(ctx context.Context, productID, adminID int64, reason, ipAddress string) error {
	err := s.repo.HideProduct(ctx, productID)
	if err != nil {
		return err
	}

	// Create audit log
	auditParams := model.CreateAuditLogParams{
		AdminID:    adminID,
		ActionType: model.ActionProductHidden,
		EntityType: model.EntityProduct,
		EntityID:   productID,
		Details: map[string]interface{}{
			"reason": reason,
		},
		IPAddress: ipAddress,
	}

	return s.repo.CreateAuditLog(ctx, auditParams)
}

// ==================== Audit Logs ====================

// GetAuditLogs retrieves audit logs
func (s *AdminService) GetAuditLogs(ctx context.Context, params model.ListAuditLogsParams) ([]model.AuditLog, int32, error) {
	// Validate params
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	return s.repo.ListAuditLogs(ctx, params)
}

// ==================== Helper Functions ====================

func isValidGroupBy(groupBy string) bool {
	return groupBy == "day" || groupBy == "week" || groupBy == "month"
}
