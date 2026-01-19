package model

import "time"

// AdminUser represents a user from admin perspective
type AdminUser struct {
	ID          int64
	Email       string
	FirstName   string
	LastName    string
	Phone       string
	Role        string // "user" or "admin"
	IsActive    bool
	IsBanned    bool
	BanReason   *string
	BannedAt    *time.Time
	BannedBy    *int64
	TotalOrders int32
	TotalSpent  float64
	LastLogin   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserStatistics represents user activity statistics
type UserStatistics struct {
	UserID       int64
	TotalBids    int32
	TotalAsks    int32
	TotalMatches int32
	TotalOrders  int32
	TotalSpent   float64
	TotalEarned  float64
}

// BanUserParams contains parameters for banning a user
type BanUserParams struct {
	UserID   int64
	AdminID  int64
	Reason   string
	BannedAt time.Time
}

// UnbanUserParams contains parameters for unbanning a user
type UnbanUserParams struct {
	UserID  int64
	AdminID int64
}

// UpdateUserRoleParams contains parameters for updating user role
type UpdateUserRoleParams struct {
	UserID  int64
	AdminID int64
	NewRole string
}

// DeleteUserParams contains parameters for deleting a user
type DeleteUserParams struct {
	UserID  int64
	AdminID int64
	Reason  string
}

// ListUsersParams contains parameters for listing users
type ListUsersParams struct {
	Page     int32
	PageSize int32
	Status   string // "all", "active", "banned"
	Role     string // "all", "user", "admin"
	Search   string // email or name search
}

// FullName returns user's full name
func (u *AdminUser) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}

// CanBeDeleted checks if user can be deleted
func (u *AdminUser) CanBeDeleted() bool {
	// Don't allow deleting admin users with orders
	if u.Role == "admin" && u.TotalOrders > 0 {
		return false
	}
	return true
}

// IsValidRole checks if role is valid
func IsValidRole(role string) bool {
	return role == "user" || role == "admin"
}
