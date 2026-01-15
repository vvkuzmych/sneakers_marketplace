package service

import (
	"context"
	"fmt"

	"github.com/vvkuzmych/sneakers_marketplace/internal/product/model"
	"github.com/vvkuzmych/sneakers_marketplace/internal/product/repository"
)

// ProductService handles business logic for products
type ProductService struct {
	productRepo   *repository.ProductRepository
	inventoryRepo *repository.InventoryRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo *repository.ProductRepository, inventoryRepo *repository.InventoryRepository) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	// Check if SKU already exists
	existing, err := s.productRepo.GetProductBySKU(ctx, product.SKU)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", product.SKU)
	}

	// Set default active status
	if !product.IsActive {
		product.IsActive = true
	}

	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product with images and sizes
func (s *ProductService) GetProduct(ctx context.Context, productID int64) (*model.ProductWithDetails, error) {
	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Get images
	images, err := s.productRepo.GetProductImages(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	// Get sizes
	sizes, err := s.inventoryRepo.GetSizesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}

	return &model.ProductWithDetails{
		Product: *product,
		Images:  convertImagePointers(images),
		Sizes:   convertSizePointers(sizes),
	}, nil
}

// ListProducts retrieves products with pagination and filters
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int, category, brand string, activeOnly bool) ([]*model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.productRepo.ListProducts(ctx, page, pageSize, category, brand, activeOnly)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	// Check if product exists
	existing, err := s.productRepo.GetProductByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Preserve SKU and timestamps
	product.SKU = existing.SKU
	product.CreatedAt = existing.CreatedAt

	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	return s.productRepo.DeleteProduct(ctx, productID)
}

// SearchProducts searches products by query
func (s *ProductService) SearchProducts(ctx context.Context, query string, page, pageSize int) ([]*model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.productRepo.SearchProducts(ctx, query, page, pageSize)
}

// AddProductImage adds an image to a product
func (s *ProductService) AddProductImage(ctx context.Context, image *model.ProductImage) (*model.ProductImage, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, image.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if err := s.productRepo.AddProductImage(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to add image: %w", err)
	}

	return image, nil
}

// DeleteProductImage deletes an image
func (s *ProductService) DeleteProductImage(ctx context.Context, imageID int64) error {
	return s.productRepo.DeleteProductImage(ctx, imageID)
}

// AddSize adds a size with initial inventory
func (s *ProductService) AddSize(ctx context.Context, size *model.Size) (*model.Size, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, size.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if err := s.inventoryRepo.AddSize(ctx, size); err != nil {
		return nil, fmt.Errorf("failed to add size: %w", err)
	}

	return size, nil
}

// GetAvailableSizes retrieves available sizes for a product
func (s *ProductService) GetAvailableSizes(ctx context.Context, productID int64) ([]*model.Size, error) {
	return s.inventoryRepo.GetSizesByProductID(ctx, productID)
}

// UpdateInventory updates inventory quantity for a size
func (s *ProductService) UpdateInventory(ctx context.Context, sizeID int64, quantity int, notes string) (*model.Size, error) {
	return s.inventoryRepo.UpdateInventory(ctx, sizeID, quantity, notes)
}

// ReserveInventory reserves inventory for an order
func (s *ProductService) ReserveInventory(ctx context.Context, sizeID int64, quantity int, orderID string) error {
	return s.inventoryRepo.ReserveInventory(ctx, sizeID, quantity, orderID)
}

// ReleaseInventory releases reserved inventory
func (s *ProductService) ReleaseInventory(ctx context.Context, sizeID int64, quantity int, orderID string) error {
	return s.inventoryRepo.ReleaseInventory(ctx, sizeID, quantity, orderID)
}

// Helper functions

func convertImagePointers(images []*model.ProductImage) []model.ProductImage {
	result := make([]model.ProductImage, len(images))
	for i, img := range images {
		result[i] = *img
	}
	return result
}

func convertSizePointers(sizes []*model.Size) []model.Size {
	result := make([]model.Size, len(sizes))
	for i, size := range sizes {
		result[i] = *size
	}
	return result
}
