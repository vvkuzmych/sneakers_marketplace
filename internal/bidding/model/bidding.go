package model

import (
	"time"
)

// Bid represents a buyer's offer to purchase at a specific price
type Bid struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	ProductID int64      `json:"product_id"`
	SizeID    int64      `json:"size_id"`
	Price     float64    `json:"price"`
	Quantity  int        `json:"quantity"`
	Status    string     `json:"status"` // active, matched, canceled, expired
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	MatchedAt *time.Time `json:"matched_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Ask represents a seller's offer to sell at a specific price
type Ask struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	ProductID int64      `json:"product_id"`
	SizeID    int64      `json:"size_id"`
	Price     float64    `json:"price"`
	Quantity  int        `json:"quantity"`
	Status    string     `json:"status"` // active, matched, canceled, expired
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	MatchedAt *time.Time `json:"matched_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Match represents a matched bid and ask
type Match struct {
	ID          int64      `json:"id"`
	BidID       int64      `json:"bid_id"`
	AskID       int64      `json:"ask_id"`
	BuyerID     int64      `json:"buyer_id"`
	SellerID    int64      `json:"seller_id"`
	ProductID   int64      `json:"product_id"`
	SizeID      int64      `json:"size_id"`
	Price       float64    `json:"price"`
	Quantity    int        `json:"quantity"`
	Status      string     `json:"status"` // pending, completed, failed
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// MarketPrice represents current market data for a product/size
type MarketPrice struct {
	ProductID  int64   `json:"product_id"`
	SizeID     int64   `json:"size_id"`
	HighestBid float64 `json:"highest_bid"`
	LowestAsk  float64 `json:"lowest_ask"`
	LastSale   float64 `json:"last_sale"`
	TotalBids  int64   `json:"total_bids"`
	TotalAsks  int64   `json:"total_asks"`
}

// Status constants
const (
	StatusActive    = "active"
	StatusMatched   = "matched"
	StatusCancelled = "canceled"
	StatusExpired   = "expired"
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// IsActive returns true if bid/ask is active
func (b *Bid) IsActive() bool {
	return b.Status == StatusActive
}

// IsActive returns true if ask is active
func (a *Ask) IsActive() bool {
	return a.Status == StatusActive
}

// CanMatch checks if a bid and ask can be matched
func CanMatch(bid *Bid, ask *Ask) bool {
	return bid.IsActive() &&
		ask.IsActive() &&
		bid.ProductID == ask.ProductID &&
		bid.SizeID == ask.SizeID &&
		bid.Price >= ask.Price &&
		bid.Quantity == ask.Quantity // For simplicity, match exact quantities
}
