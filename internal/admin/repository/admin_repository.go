package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/model"
)

// AdminRepository handles database operations for admin functionality
type AdminRepository struct {
	db *pgxpool.Pool
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

// ==================== User Management ====================

// ListUsers retrieves a paginated list of users
func (r *AdminRepository) ListUsers(ctx context.Context, params model.ListUsersParams) ([]model.AdminUser, int32, error) {
	// Build query with filters
	query := `
		SELECT id, email, first_name, last_name, phone, role, is_active, is_banned, 
		       ban_reason, banned_at, banned_by, total_orders, total_spent, last_login,
		       created_at, updated_at
		FROM users
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Apply status filter
	if params.Status == "active" {
		query += " AND is_active = true AND is_banned = false"
	} else if params.Status == "banned" {
		query += " AND is_banned = true"
	}

	// Apply role filter
	if params.Role != "" && params.Role != "all" {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, params.Role)
		argCount++
	}

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
			argCount, argCount, argCount)
		args = append(args, searchPattern)
		argCount++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS filtered"
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	offset := (params.Page - 1) * params.PageSize
	args = append(args, params.PageSize, offset)

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := []model.AdminUser{}
	for rows.Next() {
		var u model.AdminUser
		err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role,
			&u.IsActive, &u.IsBanned, &u.BanReason, &u.BannedAt, &u.BannedBy,
			&u.TotalOrders, &u.TotalSpent, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

// GetUser retrieves a single user by ID
func (r *AdminRepository) GetUser(ctx context.Context, userID int64) (*model.AdminUser, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, role, is_active, is_banned,
		       ban_reason, banned_at, banned_by, total_orders, total_spent, last_login,
		       created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u model.AdminUser
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Role,
		&u.IsActive, &u.IsBanned, &u.BanReason, &u.BannedAt, &u.BannedBy,
		&u.TotalOrders, &u.TotalSpent, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

// GetUserStatistics retrieves user activity statistics
func (r *AdminRepository) GetUserStatistics(ctx context.Context, userID int64) (*model.UserStatistics, error) {
	query := `
		SELECT 
			COALESCE(COUNT(DISTINCT b.id), 0) as total_bids,
			COALESCE(COUNT(DISTINCT a.id), 0) as total_asks,
			COALESCE(COUNT(DISTINCT m.id), 0) as total_matches,
			u.total_orders,
			u.total_spent,
			COALESCE(SUM(o.seller_fee + o.subtotal), 0) as total_earned
		FROM users u
		LEFT JOIN bids b ON b.user_id = u.id
		LEFT JOIN asks a ON a.user_id = u.id
		LEFT JOIN matches m ON m.buyer_id = u.id OR m.seller_id = u.id
		LEFT JOIN orders o ON o.seller_id = u.id AND o.status IN ('delivered', 'completed')
		WHERE u.id = $1
		GROUP BY u.id, u.total_orders, u.total_spent
	`

	var stats model.UserStatistics
	stats.UserID = userID
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.TotalBids, &stats.TotalAsks, &stats.TotalMatches,
		&stats.TotalOrders, &stats.TotalSpent, &stats.TotalEarned,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	return &stats, nil
}

// BanUser bans a user
func (r *AdminRepository) BanUser(ctx context.Context, params model.BanUserParams) error {
	query := `
		UPDATE users 
		SET is_banned = true, 
		    ban_reason = $1, 
		    banned_at = $2,
		    banned_by = $3,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND is_banned = false
	`

	result, err := r.db.Exec(ctx, query, params.Reason, params.BannedAt, params.AdminID, params.UserID)
	if err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already banned")
	}

	return nil
}

// UnbanUser unbans a user
func (r *AdminRepository) UnbanUser(ctx context.Context, params model.UnbanUserParams) error {
	query := `
		UPDATE users 
		SET is_banned = false,
		    ban_reason = NULL,
		    banned_at = NULL,
		    banned_by = NULL,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_banned = true
	`

	result, err := r.db.Exec(ctx, query, params.UserID)
	if err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or not banned")
	}

	return nil
}

// DeleteUser deletes a user (soft delete by deactivating)
func (r *AdminRepository) DeleteUser(ctx context.Context, params model.DeleteUserParams) error {
	query := `
		UPDATE users 
		SET is_active = false,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_active = true
	`

	result, err := r.db.Exec(ctx, query, params.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// UpdateUserRole updates a user's role
func (r *AdminRepository) UpdateUserRole(ctx context.Context, params model.UpdateUserRoleParams) error {
	query := `
		UPDATE users 
		SET role = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, params.NewRole, params.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ==================== Audit Logging ====================

// CreateAuditLog creates a new audit log entry
func (r *AdminRepository) CreateAuditLog(ctx context.Context, params model.CreateAuditLogParams) error {
	detailsJSON, err := json.Marshal(params.Details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}

	query := `
		INSERT INTO audit_logs (admin_id, action_type, entity_type, entity_id, details, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
	`

	_, err = r.db.Exec(ctx, query,
		params.AdminID, params.ActionType, params.EntityType,
		params.EntityID, detailsJSON, params.IPAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// ListAuditLogs retrieves a paginated list of audit logs
func (r *AdminRepository) ListAuditLogs(ctx context.Context, params model.ListAuditLogsParams) ([]model.AuditLog, int32, error) {
	query := `
		SELECT al.id, al.admin_id, u.email as admin_email, al.action_type, 
		       al.entity_type, al.entity_id, al.details, al.ip_address, al.created_at
		FROM audit_logs al
		JOIN users u ON u.id = al.admin_id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	// Apply filters
	if params.ActionType != "" && params.ActionType != "all" {
		query += fmt.Sprintf(" AND al.action_type = $%d", argCount)
		args = append(args, params.ActionType)
		argCount++
	}

	if params.AdminID > 0 {
		query += fmt.Sprintf(" AND al.admin_id = $%d", argCount)
		args = append(args, params.AdminID)
		argCount++
	}

	if params.DateFrom != nil {
		query += fmt.Sprintf(" AND al.created_at >= $%d", argCount)
		args = append(args, *params.DateFrom)
		argCount++
	}

	if params.DateTo != nil {
		query += fmt.Sprintf(" AND al.created_at <= $%d", argCount)
		args = append(args, *params.DateTo)
		argCount++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS filtered"
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY al.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	offset := (params.Page - 1) * params.PageSize
	args = append(args, params.PageSize, offset)

	// Execute query
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	logs := []model.AuditLog{}
	for rows.Next() {
		var log model.AuditLog
		var detailsJSON string
		err := rows.Scan(
			&log.ID, &log.AdminID, &log.AdminEmail, &log.ActionType,
			&log.EntityType, &log.EntityID, &detailsJSON, &log.IPAddress, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}

		// Parse details (ignore error if JSON is invalid)
		if details, err := model.ParseDetailsFromJSON(detailsJSON); err == nil {
			log.Details = details
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetAdminEmail retrieves admin email by ID
func (r *AdminRepository) GetAdminEmail(ctx context.Context, adminID int64) (string, error) {
	var email string
	err := r.db.QueryRow(ctx, "SELECT email FROM users WHERE id = $1", adminID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("failed to get admin email: %w", err)
	}
	return email, nil
}
