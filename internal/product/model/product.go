package model

import (
	"time"
)

// Product represents a sneaker product
type Product struct {
	ID          int64     `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ReleaseYear int       `json:"release_year"`
	RetailPrice float64   `json:"retail_price"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	ImageURL     string    `json:"image_url"`
	AltText      string    `json:"alt_text"`
	DisplayOrder int       `json:"display_order"`
	IsPrimary    bool      `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
}

// Size represents product inventory by size
type Size struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Size      string    `json:"size"`
	Quantity  int       `json:"quantity"`
	Reserved  int       `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Available returns the available (non-reserved) quantity
func (s *Size) Available() int {
	return s.Quantity - s.Reserved
}

// InventoryTransaction represents an inventory change audit log
type InventoryTransaction struct {
	ID              int64     `json:"id"`
	SizeID          int64     `json:"size_id"`
	TransactionType string    `json:"transaction_type"` // addition, sale, reservation, release
	QuantityChange  int       `json:"quantity_change"`
	QuantityBefore  int       `json:"quantity_before"`
	QuantityAfter   int       `json:"quantity_after"`
	ReferenceID     string    `json:"reference_id"` // order_id, purchase_order_id, etc.
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// ProductWithDetails represents a product with images and sizes
type ProductWithDetails struct {
	Product
	Images []ProductImage `json:"images"`
	Sizes  []Size         `json:"sizes"`
}
