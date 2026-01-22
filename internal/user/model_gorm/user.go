package model_gorm

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system (GORM version)
type User struct {
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Role         string    `gorm:"size:20;default:user" json:"role"`
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

// BeforeCreate hook - runs before inserting a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Set default values if needed
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}

// Address represents a user address (GORM version)
type Address struct {
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	City        string    `gorm:"size:100" json:"city"`
	StreetLine1 string    `gorm:"size:255" json:"street_line1"`
	StreetLine2 string    `gorm:"size:255" json:"street_line2"`
	State       string    `gorm:"size:100" json:"state"`
	PostalCode  string    `gorm:"size:20" json:"postal_code"`
	Country     string    `gorm:"size:100" json:"country"`
	AddressType string    `gorm:"size:20" json:"address_type"`
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"user_id"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
}

// TableName overrides the table name
func (Address) TableName() string {
	return "addresses"
}

// Session represents a user session (GORM version)
type Session struct {
	User             User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	ExpiresAt        time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	IPAddress        *string   `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent        *string   `gorm:"size:500" json:"user_agent,omitempty"`
	TokenHash        string    `gorm:"size:255;not null" json:"token_hash"`
	RefreshTokenHash string    `gorm:"size:255;not null" json:"refresh_token_hash"`
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"not null;index" json:"user_id"`
}

// TableName overrides the table name
func (Session) TableName() string {
	return "sessions"
}

// Scopes for common queries

// ActiveUsers returns only active users
func ActiveUsers(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", true)
}

// AdminUsers returns only admin users
func AdminUsers(db *gorm.DB) *gorm.DB {
	return db.Where("role = ?", "admin")
}

// VerifiedUsers returns only verified users
func VerifiedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("is_verified = ?", true)
}
