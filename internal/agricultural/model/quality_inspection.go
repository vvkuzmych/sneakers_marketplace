package model

import (
	"database/sql"
	"time"
)

// QualityInspection represents a quality inspection record
type QualityInspection struct {
	ID        int64         `json:"id"`
	ProductID int64         `json:"product_id"`
	OrderID   sql.NullInt64 `json:"order_id,omitempty"`

	// Inspector
	InspectorName     string  `json:"inspector_name"`
	InspectorLicense  *string `json:"inspector_license,omitempty"`
	InspectionCompany *string `json:"inspection_company,omitempty"`

	// Inspection details
	InspectionDate     time.Time `json:"inspection_date"`
	InspectionLocation *string   `json:"inspection_location,omitempty"`

	// Test Results (JSONB)
	TestResults JSONB `json:"test_results"` // {moisture: 13.2, protein: 12.8, ...}

	// Grade & Quality
	AssignedGrade *string  `json:"assigned_grade,omitempty"`
	QualityScore  *float64 `json:"quality_score,omitempty"` // 0-100

	// Pass/Fail
	Passed        bool     `json:"passed"`
	FailedReasons []string `json:"failed_reasons,omitempty"`

	// Certification
	CertificateNumber   *string    `json:"certificate_number,omitempty"`
	CertificateURL      *string    `json:"certificate_url,omitempty"`
	CertificateIssuedAt *time.Time `json:"certificate_issued_at,omitempty"`

	// Compliance
	USDAApproved    bool `json:"usda_approved"`
	OrganicVerified bool `json:"organic_verified"`

	// Timestamp
	CreatedAt time.Time `json:"created_at"`
}

// CreateQualityInspectionRequest is used for creating a new inspection
type CreateQualityInspectionRequest struct {
	ProductID          int64                  `json:"product_id" binding:"required"`
	OrderID            *int64                 `json:"order_id"`
	InspectorName      string                 `json:"inspector_name" binding:"required"`
	InspectorLicense   *string                `json:"inspector_license"`
	InspectionCompany  *string                `json:"inspection_company"`
	InspectionDate     time.Time              `json:"inspection_date" binding:"required"`
	InspectionLocation *string                `json:"inspection_location"`
	TestResults        map[string]interface{} `json:"test_results" binding:"required"`
	AssignedGrade      *string                `json:"assigned_grade"`
	QualityScore       *float64               `json:"quality_score"`
	Passed             bool                   `json:"passed"`
	FailedReasons      []string               `json:"failed_reasons"`
	USDAApproved       bool                   `json:"usda_approved"`
	OrganicVerified    bool                   `json:"organic_verified"`
}

// IssueCertificateRequest is used for issuing a certificate for an inspection
type IssueCertificateRequest struct {
	CertificateNumber string `json:"certificate_number" binding:"required"`
	CertificateURL    string `json:"certificate_url" binding:"required"`
}
