package model

import "time"

// AdminUser represents a user from admin perspective
type AdminUser struct {
	UpdatedAt   time.Time
	CreatedAt   time.Time
	BanReason   *string
	LastLogin   *time.Time
	BannedBy    *int64
	BannedAt    *time.Time
	Phone       string
	Role        string
	LastName    string
	FirstName   string
	Email       string
	ID          int64
	TotalSpent  float64
	TotalOrders int32
	IsBanned    bool
	IsActive    bool
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
	BannedAt time.Time
	Reason   string
	UserID   int64
	AdminID  int64
}

// UnbanUserParams contains parameters for unbanning a user
type UnbanUserParams struct {
	UserID  int64
	AdminID int64
}

// UpdateUserRoleParams contains parameters for updating user role
type UpdateUserRoleParams struct {
	NewRole string
	UserID  int64
	AdminID int64
}

// DeleteUserParams contains parameters for deleting a user
type DeleteUserParams struct {
	Reason  string
	UserID  int64
	AdminID int64
}

// ListUsersParams contains parameters for listing users
type ListUsersParams struct {
	Status   string
	Role     string
	Search   string
	Page     int32
	PageSize int32
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
