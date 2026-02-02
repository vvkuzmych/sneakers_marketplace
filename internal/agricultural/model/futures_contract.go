package model

import (
	"time"
)

// FuturesContract represents a forward or futures contract
type FuturesContract struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`

	// Contract details
	ContractNumber string `json:"contract_number"`
	ContractType   string `json:"contract_type"` // 'forward', 'futures'

	// Pricing
	StrikePrice   float64 `json:"strike_price"` // Agreed price per unit
	Quantity      float64 `json:"quantity"`     // Amount (tons, bushels)
	UnitOfMeasure string  `json:"unit_of_measure"`

	// Dates
	ContractDate   time.Time `json:"contract_date"`
	DeliveryDate   time.Time `json:"delivery_date"`
	ExpirationDate time.Time `json:"expiration_date"`

	// Parties
	BuyerID  int64 `json:"buyer_id"`
	SellerID int64 `json:"seller_id"`

	// Status
	Status       string     `json:"status"` // 'active', 'settled', 'expired', 'cancelled'
	SettledPrice *float64   `json:"settled_price,omitempty"`
	SettledAt    *time.Time `json:"settled_at,omitempty"`

	// Margin & Collateral
	MarginRequirement *float64 `json:"margin_requirement,omitempty"`
	MarginPosted      *float64 `json:"margin_posted,omitempty"`

	// Delivery
	DeliveryLocation *string `json:"delivery_location,omitempty"`
	DeliveryStatus   *string `json:"delivery_status,omitempty"` // 'pending', 'in_transit', 'delivered'

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFuturesContractRequest is used for creating a new futures contract
type CreateFuturesContractRequest struct {
	ProductID         int64     `json:"product_id" binding:"required"`
	ContractType      string    `json:"contract_type" binding:"required,oneof=forward futures"`
	StrikePrice       float64   `json:"strike_price" binding:"required,gt=0"`
	Quantity          float64   `json:"quantity" binding:"required,gt=0"`
	UnitOfMeasure     string    `json:"unit_of_measure" binding:"required"`
	DeliveryDate      time.Time `json:"delivery_date" binding:"required"`
	ExpirationDate    time.Time `json:"expiration_date" binding:"required"`
	BuyerID           int64     `json:"buyer_id" binding:"required"`
	SellerID          int64     `json:"seller_id" binding:"required"`
	MarginRequirement *float64  `json:"margin_requirement"`
	DeliveryLocation  *string   `json:"delivery_location"`
}

// SettleFuturesContractRequest is used for settling a contract
type SettleFuturesContractRequest struct {
	SettledPrice float64 `json:"settled_price" binding:"required,gt=0"`
}
