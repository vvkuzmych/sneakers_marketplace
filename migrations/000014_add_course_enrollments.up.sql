-- Migration: Add course_enrollments table
-- For tracking student enrollments and progress
-- Date: 2026-02-02

BEGIN;

CREATE TABLE course_enrollments (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Purchase details
    match_id BIGINT REFERENCES matches(id),   -- From bid-ask match
    price_paid DECIMAL(10,2) NOT NULL,
    
    -- Enrollment
    enrolled_at TIMESTAMP DEFAULT NOW(),
    access_expires_at TIMESTAMP,              -- For limited-time access
    
    -- Progress
    progress_percent DECIMAL(5,2) DEFAULT 0,  -- 0-100%
    lectures_completed INTEGER DEFAULT 0,
    quizzes_completed INTEGER DEFAULT 0,
    projects_completed INTEGER DEFAULT 0,
    last_accessed_at TIMESTAMP,
    
    -- Completion
    completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP,
    certificate_issued BOOLEAN DEFAULT false,
    certificate_url TEXT,
    certificate_issued_at TIMESTAMP,
    
    -- Ratings & Reviews
    rating INTEGER,                           -- 1-5 stars
    review_text TEXT,
    review_submitted_at TIMESTAMP,
    
    -- Refund
    refund_requested BOOLEAN DEFAULT false,
    refund_reason TEXT,
    refund_approved BOOLEAN,
    refund_processed_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(course_id, user_id),               -- One enrollment per user per course
    CHECK (progress_percent >= 0 AND progress_percent <= 100),
    CHECK (rating IS NULL OR (rating >= 1 AND rating <= 5))
);

-- Indexes
CREATE INDEX idx_course_enrollments_course ON course_enrollments(course_id);
CREATE INDEX idx_course_enrollments_user ON course_enrollments(user_id);
CREATE INDEX idx_course_enrollments_completed ON course_enrollments(completed);
CREATE INDEX idx_course_enrollments_access_expires ON course_enrollments(access_expires_at);

-- Comments
COMMENT ON TABLE course_enrollments IS 'Student enrollments and progress tracking';
COMMENT ON COLUMN course_enrollments.progress_percent IS 'Overall course completion percentage';

COMMIT;
