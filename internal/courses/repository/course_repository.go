package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vvkuzm/sneakers_marketplace/internal/courses/model"
)

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{db: db}
}

// CreateCourse creates a new course
func (r *CourseRepository) CreateCourse(ctx context.Context, req *model.CreateCourseRequest) (*model.Course, error) {
	var curriculumJSON []byte
	if req.Curriculum != nil {
		curriculumJSON, _ = json.Marshal(req.Curriculum)
	}

	query := `
		INSERT INTO courses (
			title, subtitle, description, instructor_id, instructor_name, instructor_title,
			category, subcategory, tags, level, prerequisites, format, duration_hours,
			base_price, min_acceptable_price, thumbnail_url, curriculum, learning_outcomes,
			max_students, start_date, language, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, 'draft')
		RETURNING id, vertical, title, subtitle, description, instructor_id, instructor_name, instructor_title,
			co_instructors, category, subcategory, tags, level, prerequisites, format, duration_hours,
			num_lectures, num_quizzes, num_projects, has_certificate, curriculum, learning_outcomes,
			promo_video_url, thumbnail_url, sample_videos, base_price, min_acceptable_price, currency,
			max_students, enrolled_students, min_students_to_run, start_date, end_date, schedule, timezone,
			enrollment_opens_at, enrollment_closes_at, early_bird_deadline, avg_rating, num_reviews, total_enrollments,
			completion_rate, job_placement_rate, certificate_type, accreditation_body, certificate_url_template,
			platform, content_access_duration, language, subtitles, status, is_featured,
			created_at, updated_at, published_at
	`

	course := &model.Course{}
	err := r.db.QueryRow(ctx, query,
		req.Title, req.Subtitle, req.Description, req.InstructorID, req.InstructorName, req.InstructorTitle,
		req.Category, req.Subcategory, req.Tags, req.Level, req.Prerequisites, req.Format, req.DurationHours,
		req.BasePrice, req.MinAcceptablePrice, req.ThumbnailURL, curriculumJSON, req.LearningOutcomes,
		req.MaxStudents, req.StartDate, req.Language,
	).Scan(
		&course.ID, &course.Vertical, &course.Title, &course.Subtitle, &course.Description,
		&course.InstructorID, &course.InstructorName, &course.InstructorTitle, &course.CoInstructors,
		&course.Category, &course.Subcategory, &course.Tags, &course.Level, &course.Prerequisites,
		&course.Format, &course.DurationHours, &course.NumLectures, &course.NumQuizzes, &course.NumProjects,
		&course.HasCertificate, &course.Curriculum, &course.LearningOutcomes, &course.PromoVideoURL,
		&course.ThumbnailURL, &course.SampleVideos, &course.BasePrice, &course.MinAcceptablePrice, &course.Currency,
		&course.MaxStudents, &course.EnrolledStudents, &course.MinStudentsToRun, &course.StartDate, &course.EndDate,
		&course.Schedule, &course.Timezone, &course.EnrollmentOpensAt, &course.EnrollmentClosesAt, &course.EarlyBirdDeadline,
		&course.AvgRating, &course.NumReviews, &course.TotalEnrollments, &course.CompletionRate, &course.JobPlacementRate,
		&course.CertificateType, &course.AccreditationBody, &course.CertificateURLTemplate, &course.Platform,
		&course.ContentAccessDuration, &course.Language, &course.Subtitles, &course.Status, &course.IsFeatured,
		&course.CreatedAt, &course.UpdatedAt, &course.PublishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	return course, nil
}

// GetCourseByID retrieves a course by ID
func (r *CourseRepository) GetCourseByID(ctx context.Context, id int64) (*model.Course, error) {
	query := `
		SELECT id, vertical, title, subtitle, description, instructor_id, instructor_name, instructor_title,
			co_instructors, category, subcategory, tags, level, prerequisites, format, duration_hours,
			num_lectures, num_quizzes, num_projects, has_certificate, curriculum, learning_outcomes,
			promo_video_url, thumbnail_url, sample_videos, base_price, min_acceptable_price, currency,
			max_students, enrolled_students, min_students_to_run, start_date, end_date, schedule, timezone,
			enrollment_opens_at, enrollment_closes_at, early_bird_deadline, avg_rating, num_reviews, total_enrollments,
			completion_rate, job_placement_rate, certificate_type, accreditation_body, certificate_url_template,
			platform, content_access_duration, language, subtitles, status, is_featured,
			created_at, updated_at, published_at
		FROM courses
		WHERE id = $1
	`

	course := &model.Course{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&course.ID, &course.Vertical, &course.Title, &course.Subtitle, &course.Description,
		&course.InstructorID, &course.InstructorName, &course.InstructorTitle, &course.CoInstructors,
		&course.Category, &course.Subcategory, &course.Tags, &course.Level, &course.Prerequisites,
		&course.Format, &course.DurationHours, &course.NumLectures, &course.NumQuizzes, &course.NumProjects,
		&course.HasCertificate, &course.Curriculum, &course.LearningOutcomes, &course.PromoVideoURL,
		&course.ThumbnailURL, &course.SampleVideos, &course.BasePrice, &course.MinAcceptablePrice, &course.Currency,
		&course.MaxStudents, &course.EnrolledStudents, &course.MinStudentsToRun, &course.StartDate, &course.EndDate,
		&course.Schedule, &course.Timezone, &course.EnrollmentOpensAt, &course.EnrollmentClosesAt, &course.EarlyBirdDeadline,
		&course.AvgRating, &course.NumReviews, &course.TotalEnrollments, &course.CompletionRate, &course.JobPlacementRate,
		&course.CertificateType, &course.AccreditationBody, &course.CertificateURLTemplate, &course.Platform,
		&course.ContentAccessDuration, &course.Language, &course.Subtitles, &course.Status, &course.IsFeatured,
		&course.CreatedAt, &course.UpdatedAt, &course.PublishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return course, nil
}

// ListCourses lists courses with filters and pagination
func (r *CourseRepository) ListCourses(ctx context.Context, category, level, format, status string, limit, offset int) ([]*model.Course, error) {
	query := `
		SELECT id, vertical, title, subtitle, description, instructor_id, instructor_name, instructor_title,
			co_instructors, category, subcategory, tags, level, prerequisites, format, duration_hours,
			num_lectures, num_quizzes, num_projects, has_certificate, curriculum, learning_outcomes,
			promo_video_url, thumbnail_url, sample_videos, base_price, min_acceptable_price, currency,
			max_students, enrolled_students, min_students_to_run, start_date, end_date, schedule, timezone,
			enrollment_opens_at, enrollment_closes_at, early_bird_deadline, avg_rating, num_reviews, total_enrollments,
			completion_rate, job_placement_rate, certificate_type, accreditation_body, certificate_url_template,
			platform, content_access_duration, language, subtitles, status, is_featured,
			created_at, updated_at, published_at
		FROM courses
		WHERE ($1 = '' OR category = $1)
		  AND ($2 = '' OR level = $2)
		  AND ($3 = '' OR format = $3)
		  AND ($4 = '' OR status = $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.Query(ctx, query, category, level, format, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list courses: %w", err)
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		course := &model.Course{}
		err := rows.Scan(
			&course.ID, &course.Vertical, &course.Title, &course.Subtitle, &course.Description,
			&course.InstructorID, &course.InstructorName, &course.InstructorTitle, &course.CoInstructors,
			&course.Category, &course.Subcategory, &course.Tags, &course.Level, &course.Prerequisites,
			&course.Format, &course.DurationHours, &course.NumLectures, &course.NumQuizzes, &course.NumProjects,
			&course.HasCertificate, &course.Curriculum, &course.LearningOutcomes, &course.PromoVideoURL,
			&course.ThumbnailURL, &course.SampleVideos, &course.BasePrice, &course.MinAcceptablePrice, &course.Currency,
			&course.MaxStudents, &course.EnrolledStudents, &course.MinStudentsToRun, &course.StartDate, &course.EndDate,
			&course.Schedule, &course.Timezone, &course.EnrollmentOpensAt, &course.EnrollmentClosesAt, &course.EarlyBirdDeadline,
			&course.AvgRating, &course.NumReviews, &course.TotalEnrollments, &course.CompletionRate, &course.JobPlacementRate,
			&course.CertificateType, &course.AccreditationBody, &course.CertificateURLTemplate, &course.Platform,
			&course.ContentAccessDuration, &course.Language, &course.Subtitles, &course.Status, &course.IsFeatured,
			&course.CreatedAt, &course.UpdatedAt, &course.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// UpdateCourse updates a course
func (r *CourseRepository) UpdateCourse(ctx context.Context, id int64, req *model.UpdateCourseRequest) error {
	var curriculumJSON []byte
	if req.Curriculum != nil {
		curriculumJSON, _ = json.Marshal(req.Curriculum)
	}

	query := `
		UPDATE courses
		SET title = COALESCE($1, title),
			description = COALESCE($2, description),
			base_price = COALESCE($3, base_price),
			thumbnail_url = COALESCE($4, thumbnail_url),
			status = COALESCE($5, status),
			curriculum = COALESCE($6, curriculum),
			learning_outcomes = COALESCE($7, learning_outcomes),
			updated_at = NOW()
		WHERE id = $8
	`

	_, err := r.db.Exec(ctx, query,
		req.Title, req.Description, req.BasePrice, req.ThumbnailURL, req.Status,
		curriculumJSON, req.LearningOutcomes, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	return nil
}

// PublishCourse publishes a course
func (r *CourseRepository) PublishCourse(ctx context.Context, id int64) error {
	query := `
		UPDATE courses
		SET status = 'published',
			published_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status = 'draft'
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to publish course: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("course not found or already published")
	}

	return nil
}

// IncrementEnrollment increments enrolled_students count
func (r *CourseRepository) IncrementEnrollment(ctx context.Context, id int64) error {
	query := `
		UPDATE courses
		SET enrolled_students = enrolled_students + 1,
			total_enrollments = total_enrollments + 1,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdateRating updates course rating
func (r *CourseRepository) UpdateRating(ctx context.Context, id int64, newRating int) error {
	query := `
		UPDATE courses
		SET avg_rating = (avg_rating * num_reviews + $1) / (num_reviews + 1),
			num_reviews = num_reviews + 1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, newRating, id)
	return err
}

// DeleteCourse deletes a course
func (r *CourseRepository) DeleteCourse(ctx context.Context, id int64) error {
	query := `DELETE FROM courses WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}
