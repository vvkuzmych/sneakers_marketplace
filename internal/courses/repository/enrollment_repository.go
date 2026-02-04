package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzmych/sneakers_marketplace/internal/courses/model"
)

type EnrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// CreateEnrollment creates a new enrollment
func (r *EnrollmentRepository) CreateEnrollment(ctx context.Context, req *model.CreateEnrollmentRequest) (*model.CourseEnrollment, error) {
	query := `
		INSERT INTO course_enrollments (
			course_id, user_id, match_id, price_paid
		) VALUES ($1, $2, $3, $4)
		RETURNING id, course_id, user_id, match_id, price_paid, enrolled_at, access_expires_at,
			progress_percent, lectures_completed, quizzes_completed, projects_completed, last_accessed_at,
			completed, completed_at, certificate_issued, certificate_url, certificate_issued_at,
			rating, review_text, review_submitted_at, refund_requested, refund_reason, refund_approved, refund_processed_at,
			created_at, updated_at
	`

	enrollment := &model.CourseEnrollment{}
	err := r.db.QueryRow(ctx, query, req.CourseID, req.UserID, req.MatchID, req.PricePaid).Scan(
		&enrollment.ID, &enrollment.CourseID, &enrollment.UserID, &enrollment.MatchID, &enrollment.PricePaid,
		&enrollment.EnrolledAt, &enrollment.AccessExpiresAt, &enrollment.ProgressPercent, &enrollment.LecturesCompleted,
		&enrollment.QuizzesCompleted, &enrollment.ProjectsCompleted, &enrollment.LastAccessedAt, &enrollment.Completed,
		&enrollment.CompletedAt, &enrollment.CertificateIssued, &enrollment.CertificateURL, &enrollment.CertificateIssuedAt,
		&enrollment.Rating, &enrollment.ReviewText, &enrollment.ReviewSubmittedAt, &enrollment.RefundRequested,
		&enrollment.RefundReason, &enrollment.RefundApproved, &enrollment.RefundProcessedAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	return enrollment, nil
}

// GetEnrollmentByID retrieves an enrollment by ID
func (r *EnrollmentRepository) GetEnrollmentByID(ctx context.Context, id int64) (*model.CourseEnrollment, error) {
	query := `
		SELECT id, course_id, user_id, match_id, price_paid, enrolled_at, access_expires_at,
			progress_percent, lectures_completed, quizzes_completed, projects_completed, last_accessed_at,
			completed, completed_at, certificate_issued, certificate_url, certificate_issued_at,
			rating, review_text, review_submitted_at, refund_requested, refund_reason, refund_approved, refund_processed_at,
			created_at, updated_at
		FROM course_enrollments
		WHERE id = $1
	`

	enrollment := &model.CourseEnrollment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&enrollment.ID, &enrollment.CourseID, &enrollment.UserID, &enrollment.MatchID, &enrollment.PricePaid,
		&enrollment.EnrolledAt, &enrollment.AccessExpiresAt, &enrollment.ProgressPercent, &enrollment.LecturesCompleted,
		&enrollment.QuizzesCompleted, &enrollment.ProjectsCompleted, &enrollment.LastAccessedAt, &enrollment.Completed,
		&enrollment.CompletedAt, &enrollment.CertificateIssued, &enrollment.CertificateURL, &enrollment.CertificateIssuedAt,
		&enrollment.Rating, &enrollment.ReviewText, &enrollment.ReviewSubmittedAt, &enrollment.RefundRequested,
		&enrollment.RefundReason, &enrollment.RefundApproved, &enrollment.RefundProcessedAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return enrollment, nil
}

// GetEnrollmentByUserAndCourse retrieves an enrollment by user and course
func (r *EnrollmentRepository) GetEnrollmentByUserAndCourse(ctx context.Context, userID, courseID int64) (*model.CourseEnrollment, error) {
	query := `
		SELECT id, course_id, user_id, match_id, price_paid, enrolled_at, access_expires_at,
			progress_percent, lectures_completed, quizzes_completed, projects_completed, last_accessed_at,
			completed, completed_at, certificate_issued, certificate_url, certificate_issued_at,
			rating, review_text, review_submitted_at, refund_requested, refund_reason, refund_approved, refund_processed_at,
			created_at, updated_at
		FROM course_enrollments
		WHERE user_id = $1 AND course_id = $2
	`

	enrollment := &model.CourseEnrollment{}
	err := r.db.QueryRow(ctx, query, userID, courseID).Scan(
		&enrollment.ID, &enrollment.CourseID, &enrollment.UserID, &enrollment.MatchID, &enrollment.PricePaid,
		&enrollment.EnrolledAt, &enrollment.AccessExpiresAt, &enrollment.ProgressPercent, &enrollment.LecturesCompleted,
		&enrollment.QuizzesCompleted, &enrollment.ProjectsCompleted, &enrollment.LastAccessedAt, &enrollment.Completed,
		&enrollment.CompletedAt, &enrollment.CertificateIssued, &enrollment.CertificateURL, &enrollment.CertificateIssuedAt,
		&enrollment.Rating, &enrollment.ReviewText, &enrollment.ReviewSubmittedAt, &enrollment.RefundRequested,
		&enrollment.RefundReason, &enrollment.RefundApproved, &enrollment.RefundProcessedAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return enrollment, nil
}

// ListEnrollmentsByUser retrieves all enrollments for a user
func (r *EnrollmentRepository) ListEnrollmentsByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.CourseEnrollment, error) {
	query := `
		SELECT id, course_id, user_id, match_id, price_paid, enrolled_at, access_expires_at,
			progress_percent, lectures_completed, quizzes_completed, projects_completed, last_accessed_at,
			completed, completed_at, certificate_issued, certificate_url, certificate_issued_at,
			rating, review_text, review_submitted_at, refund_requested, refund_reason, refund_approved, refund_processed_at,
			created_at, updated_at
		FROM course_enrollments
		WHERE user_id = $1
		ORDER BY enrolled_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*model.CourseEnrollment
	for rows.Next() {
		enrollment := &model.CourseEnrollment{}
		err := rows.Scan(
			&enrollment.ID, &enrollment.CourseID, &enrollment.UserID, &enrollment.MatchID, &enrollment.PricePaid,
			&enrollment.EnrolledAt, &enrollment.AccessExpiresAt, &enrollment.ProgressPercent, &enrollment.LecturesCompleted,
			&enrollment.QuizzesCompleted, &enrollment.ProjectsCompleted, &enrollment.LastAccessedAt, &enrollment.Completed,
			&enrollment.CompletedAt, &enrollment.CertificateIssued, &enrollment.CertificateURL, &enrollment.CertificateIssuedAt,
			&enrollment.Rating, &enrollment.ReviewText, &enrollment.ReviewSubmittedAt, &enrollment.RefundRequested,
			&enrollment.RefundReason, &enrollment.RefundApproved, &enrollment.RefundProcessedAt, &enrollment.CreatedAt, &enrollment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

// UpdateProgress updates student progress
func (r *EnrollmentRepository) UpdateProgress(ctx context.Context, id int64, req *model.UpdateProgressRequest) error {
	query := `
		UPDATE course_enrollments
		SET progress_percent = COALESCE($1, progress_percent),
			lectures_completed = COALESCE($2, lectures_completed),
			quizzes_completed = COALESCE($3, quizzes_completed),
			projects_completed = COALESCE($4, projects_completed),
			last_accessed_at = NOW(),
			updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query, req.ProgressPercent, req.LecturesCompleted, req.QuizzesCompleted, req.ProjectsCompleted, id)
	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	return nil
}

// MarkCompleted marks a course as completed
func (r *EnrollmentRepository) MarkCompleted(ctx context.Context, id int64) error {
	query := `
		UPDATE course_enrollments
		SET completed = true,
			completed_at = NOW(),
			progress_percent = 100.0,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark completed: %w", err)
	}

	return nil
}

// IssueCertificate issues a certificate for an enrollment
func (r *EnrollmentRepository) IssueCertificate(ctx context.Context, id int64, req *model.IssueCertificateRequest) error {
	query := `
		UPDATE course_enrollments
		SET certificate_issued = true,
			certificate_url = $1,
			certificate_issued_at = NOW(),
			updated_at = NOW()
		WHERE id = $2 AND completed = true
	`

	result, err := r.db.Exec(ctx, query, req.CertificateURL, id)
	if err != nil {
		return fmt.Errorf("failed to issue certificate: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("enrollment not found or not completed")
	}

	return nil
}

// SubmitReview submits a course review
func (r *EnrollmentRepository) SubmitReview(ctx context.Context, id int64, req *model.SubmitReviewRequest) error {
	query := `
		UPDATE course_enrollments
		SET rating = $1,
			review_text = $2,
			review_submitted_at = NOW(),
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, req.Rating, req.ReviewText, id)
	if err != nil {
		return fmt.Errorf("failed to submit review: %w", err)
	}

	return nil
}

// RequestRefund requests a refund
func (r *EnrollmentRepository) RequestRefund(ctx context.Context, id int64, req *model.RequestRefundRequest) error {
	query := `
		UPDATE course_enrollments
		SET refund_requested = true,
			refund_reason = $1,
			updated_at = NOW()
		WHERE id = $2 AND refund_requested = false
	`

	result, err := r.db.Exec(ctx, query, req.RefundReason, id)
	if err != nil {
		return fmt.Errorf("failed to request refund: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("enrollment not found or refund already requested")
	}

	return nil
}

// ApproveRefund approves a refund
func (r *EnrollmentRepository) ApproveRefund(ctx context.Context, id int64) error {
	query := `
		UPDATE course_enrollments
		SET refund_approved = true,
			refund_processed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND refund_requested = true AND refund_approved IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to approve refund: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("enrollment not found or not eligible for refund")
	}

	return nil
}

// DenyRefund denies a refund
func (r *EnrollmentRepository) DenyRefund(ctx context.Context, id int64) error {
	query := `
		UPDATE course_enrollments
		SET refund_approved = false,
			updated_at = NOW()
		WHERE id = $1 AND refund_requested = true AND refund_approved IS NULL
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to deny refund: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("enrollment not found or not eligible for refund")
	}

	return nil
}

// GetEnrollmentStats gets enrollment statistics for a course
func (r *EnrollmentRepository) GetEnrollmentStats(ctx context.Context, courseID int64) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_enrollments,
			AVG(progress_percent) as avg_progress,
			COUNT(CASE WHEN completed = true THEN 1 END) as completed_count,
			AVG(CASE WHEN rating IS NOT NULL THEN rating END) as avg_rating
		FROM course_enrollments
		WHERE course_id = $1
	`

	var totalEnrollments, completedCount int
	var avgProgress, avgRating *float64

	err := r.db.QueryRow(ctx, query, courseID).Scan(&totalEnrollments, &avgProgress, &completedCount, &avgRating)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_enrollments": totalEnrollments,
		"avg_progress":      avgProgress,
		"completed_count":   completedCount,
		"avg_rating":        avgRating,
	}

	return stats, nil
}
