package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/agricultural/model"
)

type QualityInspectionRepository struct {
	db *pgxpool.Pool
}

func NewQualityInspectionRepository(db *pgxpool.Pool) *QualityInspectionRepository {
	return &QualityInspectionRepository{db: db}
}

// CreateInspection creates a new quality inspection
func (r *QualityInspectionRepository) CreateInspection(ctx context.Context, req *model.CreateQualityInspectionRequest) (*model.QualityInspection, error) {
	var testResultsJSON []byte
	if req.TestResults != nil {
		testResultsJSON, _ = json.Marshal(req.TestResults)
	}

	query := `
		INSERT INTO quality_inspections (
			product_id, order_id, inspector_name, inspector_license, inspection_company,
			inspection_date, inspection_location, test_results, assigned_grade, quality_score,
			passed, failed_reasons, usda_approved, organic_verified
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, product_id, order_id, inspector_name, inspector_license, inspection_company,
			inspection_date, inspection_location, test_results, assigned_grade, quality_score,
			passed, failed_reasons, certificate_number, certificate_url, certificate_issued_at,
			usda_approved, organic_verified, created_at
	`

	inspection := &model.QualityInspection{}
	err := r.db.QueryRow(ctx, query,
		req.ProductID, req.OrderID, req.InspectorName, req.InspectorLicense, req.InspectionCompany,
		req.InspectionDate, req.InspectionLocation, testResultsJSON, req.AssignedGrade, req.QualityScore,
		req.Passed, req.FailedReasons, req.USDAApproved, req.OrganicVerified,
	).Scan(
		&inspection.ID, &inspection.ProductID, &inspection.OrderID, &inspection.InspectorName,
		&inspection.InspectorLicense, &inspection.InspectionCompany, &inspection.InspectionDate,
		&inspection.InspectionLocation, &inspection.TestResults, &inspection.AssignedGrade,
		&inspection.QualityScore, &inspection.Passed, &inspection.FailedReasons,
		&inspection.CertificateNumber, &inspection.CertificateURL, &inspection.CertificateIssuedAt,
		&inspection.USDAApproved, &inspection.OrganicVerified, &inspection.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create inspection: %w", err)
	}

	return inspection, nil
}

// GetInspectionByID retrieves an inspection by ID
func (r *QualityInspectionRepository) GetInspectionByID(ctx context.Context, id int64) (*model.QualityInspection, error) {
	query := `
		SELECT id, product_id, order_id, inspector_name, inspector_license, inspection_company,
			inspection_date, inspection_location, test_results, assigned_grade, quality_score,
			passed, failed_reasons, certificate_number, certificate_url, certificate_issued_at,
			usda_approved, organic_verified, created_at
		FROM quality_inspections
		WHERE id = $1
	`

	inspection := &model.QualityInspection{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&inspection.ID, &inspection.ProductID, &inspection.OrderID, &inspection.InspectorName,
		&inspection.InspectorLicense, &inspection.InspectionCompany, &inspection.InspectionDate,
		&inspection.InspectionLocation, &inspection.TestResults, &inspection.AssignedGrade,
		&inspection.QualityScore, &inspection.Passed, &inspection.FailedReasons,
		&inspection.CertificateNumber, &inspection.CertificateURL, &inspection.CertificateIssuedAt,
		&inspection.USDAApproved, &inspection.OrganicVerified, &inspection.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get inspection: %w", err)
	}

	return inspection, nil
}

// ListInspectionsByProduct retrieves all inspections for a product
func (r *QualityInspectionRepository) ListInspectionsByProduct(ctx context.Context, productID int64) ([]*model.QualityInspection, error) {
	query := `
		SELECT id, product_id, order_id, inspector_name, inspector_license, inspection_company,
			inspection_date, inspection_location, test_results, assigned_grade, quality_score,
			passed, failed_reasons, certificate_number, certificate_url, certificate_issued_at,
			usda_approved, organic_verified, created_at
		FROM quality_inspections
		WHERE product_id = $1
		ORDER BY inspection_date DESC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections: %w", err)
	}
	defer rows.Close()

	var inspections []*model.QualityInspection
	for rows.Next() {
		inspection := &model.QualityInspection{}
		err := rows.Scan(
			&inspection.ID, &inspection.ProductID, &inspection.OrderID, &inspection.InspectorName,
			&inspection.InspectorLicense, &inspection.InspectionCompany, &inspection.InspectionDate,
			&inspection.InspectionLocation, &inspection.TestResults, &inspection.AssignedGrade,
			&inspection.QualityScore, &inspection.Passed, &inspection.FailedReasons,
			&inspection.CertificateNumber, &inspection.CertificateURL, &inspection.CertificateIssuedAt,
			&inspection.USDAApproved, &inspection.OrganicVerified, &inspection.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inspection: %w", err)
		}
		inspections = append(inspections, inspection)
	}

	return inspections, nil
}

// IssueCertificate issues a certificate for an inspection
func (r *QualityInspectionRepository) IssueCertificate(ctx context.Context, id int64, certificateNumber, certificateURL string) error {
	query := `
		UPDATE quality_inspections
		SET certificate_number = $1,
			certificate_url = $2,
			certificate_issued_at = NOW()
		WHERE id = $3 AND passed = true
	`

	result, err := r.db.Exec(ctx, query, certificateNumber, certificateURL, id)
	if err != nil {
		return fmt.Errorf("failed to issue certificate: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("inspection not found or did not pass")
	}

	return nil
}

// GetLatestInspectionByProduct gets the most recent inspection for a product
func (r *QualityInspectionRepository) GetLatestInspectionByProduct(ctx context.Context, productID int64) (*model.QualityInspection, error) {
	query := `
		SELECT id, product_id, order_id, inspector_name, inspector_license, inspection_company,
			inspection_date, inspection_location, test_results, assigned_grade, quality_score,
			passed, failed_reasons, certificate_number, certificate_url, certificate_issued_at,
			usda_approved, organic_verified, created_at
		FROM quality_inspections
		WHERE product_id = $1
		ORDER BY inspection_date DESC
		LIMIT 1
	`

	inspection := &model.QualityInspection{}
	err := r.db.QueryRow(ctx, query, productID).Scan(
		&inspection.ID, &inspection.ProductID, &inspection.OrderID, &inspection.InspectorName,
		&inspection.InspectorLicense, &inspection.InspectionCompany, &inspection.InspectionDate,
		&inspection.InspectionLocation, &inspection.TestResults, &inspection.AssignedGrade,
		&inspection.QualityScore, &inspection.Passed, &inspection.FailedReasons,
		&inspection.CertificateNumber, &inspection.CertificateURL, &inspection.CertificateIssuedAt,
		&inspection.USDAApproved, &inspection.OrganicVerified, &inspection.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get latest inspection: %w", err)
	}

	return inspection, nil
}
