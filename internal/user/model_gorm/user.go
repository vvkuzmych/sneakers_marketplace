package model_gorm

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system (GORM version)
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Role         string    `gorm:"size:20;default:user" json:"role"` // user, admin
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"` // Relation
	AddressType string    `gorm:"size:20" json:"address_type"`                            // shipping, billing
	StreetLine1 string    `gorm:"size:255" json:"street_line1"`
	StreetLine2 string    `gorm:"size:255" json:"street_line2"`
	City        string    `gorm:"size:100" json:"city"`
	State       string    `gorm:"size:100" json:"state"`
	PostalCode  string    `gorm:"size:20" json:"postal_code"`
	Country     string    `gorm:"size:100" json:"country"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName overrides the table name
func (Address) TableName() string {
	return "addresses"
}

// Session represents a user session (GORM version)
type Session struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64     `gorm:"not null;index" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	TokenHash        string    `gorm:"size:255;not null" json:"token_hash"`
	RefreshTokenHash string    `gorm:"size:255;not null" json:"refresh_token_hash"`
	IPAddress        *string   `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent        *string   `gorm:"size:500" json:"user_agent,omitempty"`
	ExpiresAt        time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
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
