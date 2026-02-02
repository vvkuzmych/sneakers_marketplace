package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Course represents an educational course
type Course struct {
	ID       int64  `json:"id"`
	Vertical string `json:"vertical"` // "education"

	// Basic Info
	Title       string  `json:"title"`
	Subtitle    *string `json:"subtitle,omitempty"`
	Description string  `json:"description"`

	// Instructor
	InstructorID    int64   `json:"instructor_id"`
	InstructorName  *string `json:"instructor_name,omitempty"`
	InstructorTitle *string `json:"instructor_title,omitempty"`
	CoInstructors   []int64 `json:"co_instructors,omitempty"`

	// Category & Tags
	Category    string   `json:"category"` // 'technology', 'business', 'design'
	Subcategory *string  `json:"subcategory,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	// Difficulty & Prerequisites
	Level         string   `json:"level"` // 'beginner', 'intermediate', 'advanced'
	Prerequisites []string `json:"prerequisites,omitempty"`

	// Course Structure
	Format         string   `json:"format"` // 'self-paced', 'cohort-based', 'live'
	DurationHours  *float64 `json:"duration_hours,omitempty"`
	NumLectures    *int     `json:"num_lectures,omitempty"`
	NumQuizzes     *int     `json:"num_quizzes,omitempty"`
	NumProjects    *int     `json:"num_projects,omitempty"`
	HasCertificate bool     `json:"has_certificate"`

	// Curriculum (JSONB)
	Curriculum       JSONB    `json:"curriculum,omitempty"`
	LearningOutcomes []string `json:"learning_outcomes,omitempty"`

	// Media
	PromoVideoURL *string `json:"promo_video_url,omitempty"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	SampleVideos  JSONB   `json:"sample_videos,omitempty"`

	// Pricing & Seats
	BasePrice          float64  `json:"base_price"`
	MinAcceptablePrice *float64 `json:"min_acceptable_price,omitempty"`
	Currency           string   `json:"currency"`

	// Seats
	MaxStudents      *int `json:"max_students,omitempty"` // NULL = unlimited
	EnrolledStudents int  `json:"enrolled_students"`
	MinStudentsToRun *int `json:"min_students_to_run,omitempty"`

	// Schedule
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Schedule  *string    `json:"schedule,omitempty"`
	Timezone  *string    `json:"timezone,omitempty"`

	// Enrollment window
	EnrollmentOpensAt  *time.Time `json:"enrollment_opens_at,omitempty"`
	EnrollmentClosesAt *time.Time `json:"enrollment_closes_at,omitempty"`
	EarlyBirdDeadline  *time.Time `json:"early_bird_deadline,omitempty"`

	// Ratings
	AvgRating        float64 `json:"avg_rating"`
	NumReviews       int     `json:"num_reviews"`
	TotalEnrollments int     `json:"total_enrollments"`

	// Completion & Success
	CompletionRate   *float64 `json:"completion_rate,omitempty"`
	JobPlacementRate *float64 `json:"job_placement_rate,omitempty"`

	// Certification
	CertificateType        *string `json:"certificate_type,omitempty"`
	AccreditationBody      *string `json:"accreditation_body,omitempty"`
	CertificateURLTemplate *string `json:"certificate_url_template,omitempty"`

	// Content delivery
	Platform              *string `json:"platform,omitempty"`
	ContentAccessDuration *int    `json:"content_access_duration,omitempty"` // Days

	// Language
	Language  string   `json:"language"`
	Subtitles []string `json:"subtitles,omitempty"`

	// Status
	Status     string `json:"status"` // 'draft', 'published', 'archived'
	IsFeatured bool   `json:"is_featured"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
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

// CreateCourseRequest is used for creating a new course
type CreateCourseRequest struct {
	Title              string                 `json:"title" binding:"required"`
	Subtitle           *string                `json:"subtitle"`
	Description        string                 `json:"description" binding:"required"`
	InstructorID       int64                  `json:"instructor_id" binding:"required"`
	InstructorName     *string                `json:"instructor_name"`
	InstructorTitle    *string                `json:"instructor_title"`
	Category           string                 `json:"category" binding:"required"`
	Subcategory        *string                `json:"subcategory"`
	Tags               []string               `json:"tags"`
	Level              string                 `json:"level" binding:"required,oneof=beginner intermediate advanced"`
	Prerequisites      []string               `json:"prerequisites"`
	Format             string                 `json:"format" binding:"required,oneof=self-paced cohort-based live hybrid"`
	DurationHours      *float64               `json:"duration_hours"`
	BasePrice          float64                `json:"base_price" binding:"required,gt=0"`
	MinAcceptablePrice *float64               `json:"min_acceptable_price"`
	ThumbnailURL       string                 `json:"thumbnail_url" binding:"required"`
	Curriculum         map[string]interface{} `json:"curriculum"`
	LearningOutcomes   []string               `json:"learning_outcomes"`
	MaxStudents        *int                   `json:"max_students"`
	StartDate          *time.Time             `json:"start_date"`
	Language           string                 `json:"language"`
}

// UpdateCourseRequest is used for updating a course
type UpdateCourseRequest struct {
	Title            *string                `json:"title"`
	Description      *string                `json:"description"`
	BasePrice        *float64               `json:"base_price"`
	ThumbnailURL     *string                `json:"thumbnail_url"`
	Status           *string                `json:"status" binding:"omitempty,oneof=draft published archived"`
	Curriculum       map[string]interface{} `json:"curriculum"`
	LearningOutcomes []string               `json:"learning_outcomes"`
}
