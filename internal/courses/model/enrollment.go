package model

import (
	"database/sql"
	"time"
)

// CourseEnrollment represents a student's enrollment in a course
type CourseEnrollment struct {
	ID       int64 `json:"id"`
	CourseID int64 `json:"course_id"`
	UserID   int64 `json:"user_id"`

	// Purchase details
	MatchID   sql.NullInt64 `json:"match_id,omitempty"`
	PricePaid float64       `json:"price_paid"`

	// Enrollment
	EnrolledAt      time.Time  `json:"enrolled_at"`
	AccessExpiresAt *time.Time `json:"access_expires_at,omitempty"`

	// Progress
	ProgressPercent   float64    `json:"progress_percent"`
	LecturesCompleted int        `json:"lectures_completed"`
	QuizzesCompleted  int        `json:"quizzes_completed"`
	ProjectsCompleted int        `json:"projects_completed"`
	LastAccessedAt    *time.Time `json:"last_accessed_at,omitempty"`

	// Completion
	Completed           bool       `json:"completed"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	CertificateIssued   bool       `json:"certificate_issued"`
	CertificateURL      *string    `json:"certificate_url,omitempty"`
	CertificateIssuedAt *time.Time `json:"certificate_issued_at,omitempty"`

	// Ratings & Reviews
	Rating            *int       `json:"rating,omitempty"` // 1-5 stars
	ReviewText        *string    `json:"review_text,omitempty"`
	ReviewSubmittedAt *time.Time `json:"review_submitted_at,omitempty"`

	// Refund
	RefundRequested   bool       `json:"refund_requested"`
	RefundReason      *string    `json:"refund_reason,omitempty"`
	RefundApproved    *bool      `json:"refund_approved,omitempty"`
	RefundProcessedAt *time.Time `json:"refund_processed_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateEnrollmentRequest is used for enrolling a student
type CreateEnrollmentRequest struct {
	CourseID  int64   `json:"course_id" binding:"required"`
	UserID    int64   `json:"user_id" binding:"required"`
	MatchID   *int64  `json:"match_id"`
	PricePaid float64 `json:"price_paid" binding:"required,gt=0"`
}

// UpdateProgressRequest is used for updating student progress
type UpdateProgressRequest struct {
	ProgressPercent   *float64 `json:"progress_percent"`
	LecturesCompleted *int     `json:"lectures_completed"`
	QuizzesCompleted  *int     `json:"quizzes_completed"`
	ProjectsCompleted *int     `json:"projects_completed"`
}

// SubmitReviewRequest is used for submitting a course review
type SubmitReviewRequest struct {
	Rating     int     `json:"rating" binding:"required,min=1,max=5"`
	ReviewText *string `json:"review_text"`
}

// RequestRefundRequest is used for requesting a refund
type RequestRefundRequest struct {
	RefundReason string `json:"refund_reason" binding:"required"`
}

// IssueCertificateRequest is used for issuing a certificate
type IssueCertificateRequest struct {
	CertificateURL string `json:"certificate_url" binding:"required"`
}
