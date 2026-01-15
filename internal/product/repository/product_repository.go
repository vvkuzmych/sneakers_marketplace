package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/product/model"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (sku, name, brand, model, color, description, category, release_year, retail_price, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		product.SKU,
		product.Name,
		product.Brand,
		product.Model,
		product.Color,
		product.Description,
		product.Category,
		product.ReleaseYear,
		product.RetailPrice,
		product.IsActive,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetProductByID retrieves a product by ID
func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, sku, name, brand, model, color, description, category, 
		       release_year, retail_price, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &model.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Brand,
		&product.Model,
		&product.Color,
		&product.Description,
		&product.Category,
		&product.ReleaseYear,
		&product.RetailPrice,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetProductBySKU retrieves a product by SKU
func (r *ProductRepository) GetProductBySKU(ctx context.Context, sku string) (*model.Product, error) {
	query := `
		SELECT id, sku, name, brand, model, color, description, category, 
		       release_year, retail_price, is_active, created_at, updated_at
		FROM products
		WHERE sku = $1
	`

	product := &model.Product{}
	err := r.db.QueryRow(ctx, query, sku).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Brand,
		&product.Model,
		&product.Color,
		&product.Description,
		&product.Category,
		&product.ReleaseYear,
		&product.RetailPrice,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

// ListProducts retrieves products with pagination and filters
func (r *ProductRepository) ListProducts(ctx context.Context, page, pageSize int, category, brand string, activeOnly bool) ([]*model.Product, int64, error) {
	// Build query dynamically based on filters
	query := `
		SELECT id, sku, name, brand, model, color, description, category, 
		       release_year, retail_price, is_active, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`

	args := []interface{}{}
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if brand != "" {
		query += fmt.Sprintf(" AND brand = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND brand = $%d", argIdx)
		args = append(args, brand)
		argIdx++
	}

	if activeOnly {
		query += " AND is_active = true"
		countQuery += " AND is_active = true"
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		product := &model.Product{}
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Brand,
			&product.Model,
			&product.Color,
			&product.Description,
			&product.Category,
			&product.ReleaseYear,
			&product.RetailPrice,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, total, nil
}

// UpdateProduct updates a product
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET name = $1, brand = $2, model = $3, color = $4, description = $5,
		    category = $6, release_year = $7, retail_price = $8, is_active = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Brand,
		product.Model,
		product.Color,
		product.Description,
		product.Category,
		product.ReleaseYear,
		product.RetailPrice,
		product.IsActive,
		product.ID,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product
func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// SearchProducts searches products by name, brand, or model
func (r *ProductRepository) SearchProducts(ctx context.Context, searchQuery string, page, pageSize int) ([]*model.Product, int64, error) {
	query := `
		SELECT id, sku, name, brand, model, color, description, category, 
		       release_year, retail_price, is_active, created_at, updated_at
		FROM products
		WHERE (name ILIKE $1 OR brand ILIKE $1 OR model ILIKE $1)
		  AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE (name ILIKE $1 OR brand ILIKE $1 OR model ILIKE $1)
		  AND is_active = true
	`

	searchPattern := "%" + searchQuery + "%"

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get products
	rows, err := r.db.Query(ctx, query, searchPattern, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		product := &model.Product{}
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Brand,
			&product.Model,
			&product.Color,
			&product.Description,
			&product.Category,
			&product.ReleaseYear,
			&product.RetailPrice,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, total, nil
}

// AddProductImage adds a new image to a product
func (r *ProductRepository) AddProductImage(ctx context.Context, image *model.ProductImage) error {
	query := `
		INSERT INTO product_images (product_id, image_url, alt_text, display_order, is_primary)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		image.ProductID,
		image.ImageURL,
		image.AltText,
		image.DisplayOrder,
		image.IsPrimary,
	).Scan(&image.ID, &image.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add product image: %w", err)
	}

	// If this is primary, unset other primary images
	if image.IsPrimary {
		if err := r.unsetOtherPrimaryImages(ctx, image.ProductID, image.ID); err != nil {
			return err
		}
	}

	return nil
}

// GetProductImages retrieves all images for a product
func (r *ProductRepository) GetProductImages(ctx context.Context, productID int64) ([]*model.ProductImage, error) {
	query := `
		SELECT id, product_id, image_url, alt_text, display_order, is_primary, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}
	defer rows.Close()

	var images []*model.ProductImage
	for rows.Next() {
		image := &model.ProductImage{}
		err := rows.Scan(
			&image.ID,
			&image.ProductID,
			&image.ImageURL,
			&image.AltText,
			&image.DisplayOrder,
			&image.IsPrimary,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		images = append(images, image)
	}

	return images, nil
}

// DeleteProductImage deletes an image
func (r *ProductRepository) DeleteProductImage(ctx context.Context, imageID int64) error {
	query := `DELETE FROM product_images WHERE id = $1`

	result, err := r.db.Exec(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// Helper: unset other primary images
func (r *ProductRepository) unsetOtherPrimaryImages(ctx context.Context, productID, imageID int64) error {
	query := `
		UPDATE product_images
		SET is_primary = false
		WHERE product_id = $1 AND id != $2 AND is_primary = true
	`

	_, err := r.db.Exec(ctx, query, productID, imageID)
	if err != nil {
		return fmt.Errorf("failed to unset other primary images: %w", err)
	}

	return nil
}
