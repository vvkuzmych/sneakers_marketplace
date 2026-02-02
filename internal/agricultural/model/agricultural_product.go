package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AgriculturalProduct represents an agricultural commodity
type AgriculturalProduct struct {
	ID       int64  `json:"id"`
	Vertical string `json:"vertical"` // "agriculture"

	// Basic Info
	CommodityType string  `json:"commodity_type"` // 'wheat', 'corn', 'soybeans'
	Variety       *string `json:"variety,omitempty"`
	Grade         *string `json:"grade,omitempty"`

	// Quantity & Units
	UnitOfMeasure    string   `json:"unit_of_measure"` // 'tons', 'bushels', 'kg'
	MinOrderQuantity *float64 `json:"min_order_quantity,omitempty"`

	// Origin
	CountryOfOrigin *string `json:"country_of_origin,omitempty"`
	StateProvince   *string `json:"state_province,omitempty"`
	FarmName        *string `json:"farm_name,omitempty"`

	// Certifications & Quality (JSONB)
	Certifications JSONB `json:"certifications,omitempty"` // ["organic", "non-gmo"]
	QualitySpecs   JSONB `json:"quality_specs,omitempty"`  // {protein: 12.5%, moisture: 13%}

	// Lab Testing
	LabTested         bool       `json:"lab_tested"`
	LabCertificateURL *string    `json:"lab_certificate_url,omitempty"`
	TestDate          *time.Time `json:"test_date,omitempty"`

	// Harvest
	HarvestYear   *int    `json:"harvest_year,omitempty"`
	HarvestSeason *string `json:"harvest_season,omitempty"`

	// Storage
	StorageLocation *string `json:"storage_location,omitempty"`
	StorageType     *string `json:"storage_type,omitempty"` // 'silo', 'warehouse'

	// Compliance
	USDACertified    bool `json:"usda_certified"`
	OrganicCertified bool `json:"organic_certified"`
	NonGMOCertified  bool `json:"non_gmo_certified"`

	// Images
	Images JSONB `json:"images,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JSONB type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value implements driver.Valuer for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	*j = data
	return nil
}

// CreateAgriculturalProductRequest is used for creating a new product
type CreateAgriculturalProductRequest struct {
	CommodityType    string                 `json:"commodity_type" binding:"required"`
	Variety          *string                `json:"variety"`
	Grade            *string                `json:"grade"`
	UnitOfMeasure    string                 `json:"unit_of_measure" binding:"required"`
	MinOrderQuantity *float64               `json:"min_order_quantity"`
	CountryOfOrigin  *string                `json:"country_of_origin"`
	StateProvince    *string                `json:"state_province"`
	FarmName         *string                `json:"farm_name"`
	Certifications   map[string]interface{} `json:"certifications"`
	QualitySpecs     map[string]interface{} `json:"quality_specs"`
	HarvestYear      *int                   `json:"harvest_year"`
	HarvestSeason    *string                `json:"harvest_season"`
	StorageLocation  *string                `json:"storage_location"`
	StorageType      *string                `json:"storage_type"`
	Images           map[string]interface{} `json:"images"`
}

// UpdateAgriculturalProductRequest is used for updating a product
type UpdateAgriculturalProductRequest struct {
	Variety           *string                `json:"variety"`
	Grade             *string                `json:"grade"`
	MinOrderQuantity  *float64               `json:"min_order_quantity"`
	QualitySpecs      map[string]interface{} `json:"quality_specs"`
	StorageLocation   *string                `json:"storage_location"`
	StorageType       *string                `json:"storage_type"`
	LabTested         *bool                  `json:"lab_tested"`
	LabCertificateURL *string                `json:"lab_certificate_url"`
	Images            map[string]interface{} `json:"images"`
}
