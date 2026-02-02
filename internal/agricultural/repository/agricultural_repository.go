package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/agricultural/model"
)

type AgriculturalRepository struct {
	db *pgxpool.Pool
}

func NewAgriculturalRepository(db *pgxpool.Pool) *AgriculturalRepository {
	return &AgriculturalRepository{db: db}
}

// CreateProduct creates a new agricultural product
func (r *AgriculturalRepository) CreateProduct(ctx context.Context, req *model.CreateAgriculturalProductRequest) (*model.AgriculturalProduct, error) {
	query := `
		INSERT INTO agricultural_products (
			commodity_type, variety, grade, unit_of_measure, min_order_quantity,
			country_of_origin, state_province, farm_name, certifications, quality_specs,
			harvest_year, harvest_season, storage_location, storage_type, images
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, vertical, commodity_type, variety, grade, unit_of_measure, min_order_quantity,
			country_of_origin, state_province, farm_name, certifications, quality_specs,
			lab_tested, lab_certificate_url, test_date, harvest_year, harvest_season,
			storage_location, storage_type, usda_certified, organic_certified, non_gmo_certified,
			images, created_at, updated_at
	`

	var certificationsJSON, qualitySpecsJSON, imagesJSON []byte
	if req.Certifications != nil {
		certificationsJSON, _ = json.Marshal(req.Certifications)
	}
	if req.QualitySpecs != nil {
		qualitySpecsJSON, _ = json.Marshal(req.QualitySpecs)
	}
	if req.Images != nil {
		imagesJSON, _ = json.Marshal(req.Images)
	}

	product := &model.AgriculturalProduct{}
	err := r.db.QueryRow(ctx, query,
		req.CommodityType, req.Variety, req.Grade, req.UnitOfMeasure, req.MinOrderQuantity,
		req.CountryOfOrigin, req.StateProvince, req.FarmName, certificationsJSON, qualitySpecsJSON,
		req.HarvestYear, req.HarvestSeason, req.StorageLocation, req.StorageType, imagesJSON,
	).Scan(
		&product.ID, &product.Vertical, &product.CommodityType, &product.Variety, &product.Grade,
		&product.UnitOfMeasure, &product.MinOrderQuantity, &product.CountryOfOrigin, &product.StateProvince,
		&product.FarmName, &product.Certifications, &product.QualitySpecs, &product.LabTested,
		&product.LabCertificateURL, &product.TestDate, &product.HarvestYear, &product.HarvestSeason,
		&product.StorageLocation, &product.StorageType, &product.USDACertified, &product.OrganicCertified,
		&product.NonGMOCertified, &product.Images, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create agricultural product: %w", err)
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (r *AgriculturalRepository) GetProductByID(ctx context.Context, id int64) (*model.AgriculturalProduct, error) {
	query := `
		SELECT id, vertical, commodity_type, variety, grade, unit_of_measure, min_order_quantity,
			country_of_origin, state_province, farm_name, certifications, quality_specs,
			lab_tested, lab_certificate_url, test_date, harvest_year, harvest_season,
			storage_location, storage_type, usda_certified, organic_certified, non_gmo_certified,
			images, created_at, updated_at
		FROM agricultural_products
		WHERE id = $1
	`

	product := &model.AgriculturalProduct{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.Vertical, &product.CommodityType, &product.Variety, &product.Grade,
		&product.UnitOfMeasure, &product.MinOrderQuantity, &product.CountryOfOrigin, &product.StateProvince,
		&product.FarmName, &product.Certifications, &product.QualitySpecs, &product.LabTested,
		&product.LabCertificateURL, &product.TestDate, &product.HarvestYear, &product.HarvestSeason,
		&product.StorageLocation, &product.StorageType, &product.USDACertified, &product.OrganicCertified,
		&product.NonGMOCertified, &product.Images, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// ListProducts lists all products with pagination
func (r *AgriculturalRepository) ListProducts(ctx context.Context, commodityType string, limit, offset int) ([]*model.AgriculturalProduct, error) {
	query := `
		SELECT id, vertical, commodity_type, variety, grade, unit_of_measure, min_order_quantity,
			country_of_origin, state_province, farm_name, certifications, quality_specs,
			lab_tested, lab_certificate_url, test_date, harvest_year, harvest_season,
			storage_location, storage_type, usda_certified, organic_certified, non_gmo_certified,
			images, created_at, updated_at
		FROM agricultural_products
		WHERE ($1 = '' OR commodity_type = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, commodityType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*model.AgriculturalProduct
	for rows.Next() {
		product := &model.AgriculturalProduct{}
		err := rows.Scan(
			&product.ID, &product.Vertical, &product.CommodityType, &product.Variety, &product.Grade,
			&product.UnitOfMeasure, &product.MinOrderQuantity, &product.CountryOfOrigin, &product.StateProvince,
			&product.FarmName, &product.Certifications, &product.QualitySpecs, &product.LabTested,
			&product.LabCertificateURL, &product.TestDate, &product.HarvestYear, &product.HarvestSeason,
			&product.StorageLocation, &product.StorageType, &product.USDACertified, &product.OrganicCertified,
			&product.NonGMOCertified, &product.Images, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (r *AgriculturalRepository) UpdateProduct(ctx context.Context, id int64, req *model.UpdateAgriculturalProductRequest) error {
	query := `
		UPDATE agricultural_products
		SET variety = COALESCE($1, variety),
			grade = COALESCE($2, grade),
			min_order_quantity = COALESCE($3, min_order_quantity),
			quality_specs = COALESCE($4, quality_specs),
			storage_location = COALESCE($5, storage_location),
			storage_type = COALESCE($6, storage_type),
			lab_tested = COALESCE($7, lab_tested),
			lab_certificate_url = COALESCE($8, lab_certificate_url),
			images = COALESCE($9, images),
			updated_at = NOW()
		WHERE id = $10
	`

	var qualitySpecsJSON, imagesJSON []byte
	if req.QualitySpecs != nil {
		qualitySpecsJSON, _ = json.Marshal(req.QualitySpecs)
	}
	if req.Images != nil {
		imagesJSON, _ = json.Marshal(req.Images)
	}

	_, err := r.db.Exec(ctx, query,
		req.Variety, req.Grade, req.MinOrderQuantity, qualitySpecsJSON,
		req.StorageLocation, req.StorageType, req.LabTested, req.LabCertificateURL,
		imagesJSON, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product by ID
func (r *AgriculturalRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM agricultural_products WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
