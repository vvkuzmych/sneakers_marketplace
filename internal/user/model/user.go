package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	ID           int64     `json:"id"`
	IsVerified   bool      `json:"is_verified"`
	IsActive     bool      `json:"is_active"`
}

// Address represents a user address
type Address struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AddressType string    `json:"address_type"`
	StreetLine1 string    `json:"street_line1"`
	StreetLine2 string    `json:"street_line2"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	IsDefault   bool      `json:"is_default"`
}

// Session represents a user session (for JWT token management)
type Session struct {
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
	IPAddress        *string   `json:"ip_address,omitempty"`
	UserAgent        *string   `json:"user_agent,omitempty"`
	TokenHash        string    `json:"token_hash"`
	RefreshTokenHash string    `json:"refresh_token_hash"`
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
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
