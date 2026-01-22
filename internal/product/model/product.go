package model

import (
	"time"
)

// Product represents a sneaker product
type Product struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	ID          int64     `json:"id"`
	ReleaseYear int       `json:"release_year"`
	RetailPrice float64   `json:"retail_price"`
	IsActive    bool      `json:"is_active"`
}

// ProductImage represents a product image
type ProductImage struct {
	CreatedAt    time.Time `json:"created_at"`
	ImageURL     string    `json:"image_url"`
	AltText      string    `json:"alt_text"`
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	DisplayOrder int       `json:"display_order"`
	IsPrimary    bool      `json:"is_primary"`
}

// Size represents product inventory by size
type Size struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Size      string    `json:"size"`
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
}

// Available returns the available (non-reserved) quantity
func (s *Size) Available() int {
	return s.Quantity - s.Reserved
}

// InventoryTransaction represents an inventory change audit log
type InventoryTransaction struct {
	CreatedAt       time.Time `json:"created_at"`
	TransactionType string    `json:"transaction_type"`
	ReferenceID     string    `json:"reference_id"`
	Notes           string    `json:"notes"`
	ID              int64     `json:"id"`
	SizeID          int64     `json:"size_id"`
	QuantityChange  int       `json:"quantity_change"`
	QuantityBefore  int       `json:"quantity_before"`
	QuantityAfter   int       `json:"quantity_after"`
}

// ProductWithDetails represents a product with images and sizes
type ProductWithDetails struct {
	Product
	Images []ProductImage `json:"images"`
	Sizes  []Size         `json:"sizes"`
}
