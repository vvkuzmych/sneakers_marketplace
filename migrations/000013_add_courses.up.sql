-- Migration: Add courses table
-- Vertical: Education / Course Marketplace
-- Date: 2026-02-02

BEGIN;

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    vertical VARCHAR(50) DEFAULT 'education',
    
    -- Basic Info
    title VARCHAR(500) NOT NULL,
    subtitle TEXT,
    description TEXT NOT NULL,
    
    -- Instructor
    instructor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    instructor_name VARCHAR(255),
    instructor_title VARCHAR(255),
    co_instructors BIGINT[],                  -- Multiple instructors
    
    -- Category & Tags
    category VARCHAR(100) NOT NULL,           -- 'technology', 'business', 'design'
    subcategory VARCHAR(100),
    tags TEXT[],
    
    -- Difficulty & Prerequisites
    level VARCHAR(50) NOT NULL,               -- 'beginner', 'intermediate', 'advanced'
    prerequisites TEXT[],
    
    -- Course Structure
    format VARCHAR(50) NOT NULL,              -- 'self-paced', 'cohort-based', 'live'
    duration_hours DECIMAL(5,1),
    num_lectures INTEGER,
    num_quizzes INTEGER,
    num_projects INTEGER,
    has_certificate BOOLEAN DEFAULT true,
    
    -- Curriculum (JSONB for flexibility)
    curriculum JSONB,
    learning_outcomes TEXT[],
    
    -- Media
    promo_video_url TEXT,
    thumbnail_url TEXT NOT NULL,
    sample_videos JSONB,
    
    -- Pricing & Seats
    base_price DECIMAL(10,2) NOT NULL,
    min_acceptable_price DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Seats (for cohort-based)
    max_students INTEGER,                     -- NULL = unlimited (self-paced)
    enrolled_students INTEGER DEFAULT 0,
    min_students_to_run INTEGER,
    
    -- Schedule
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    schedule TEXT,
    timezone VARCHAR(50),
    
    -- Enrollment window
    enrollment_opens_at TIMESTAMP,
    enrollment_closes_at TIMESTAMP,
    early_bird_deadline TIMESTAMP,
    
    -- Ratings
    avg_rating DECIMAL(3,2) DEFAULT 0.00,
    num_reviews INTEGER DEFAULT 0,
    total_enrollments INTEGER DEFAULT 0,
    
    -- Completion & Success
    completion_rate DECIMAL(5,2),
    job_placement_rate DECIMAL(5,2),
    
    -- Certification
    certificate_type VARCHAR(100),
    accreditation_body VARCHAR(255),
    certificate_url_template TEXT,
    
    -- Content delivery
    platform VARCHAR(50),
    content_access_duration INTEGER,          -- Days
    
    -- Language
    language VARCHAR(50) DEFAULT 'en',
    subtitles TEXT[],
    
    -- Status
    status VARCHAR(50) NOT NULL,              -- 'draft', 'published', 'archived'
    is_featured BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    published_at TIMESTAMP,
    
    CHECK (base_price > 0),
    CHECK (enrolled_students <= max_students OR max_students IS NULL),
    CHECK (avg_rating >= 0 AND avg_rating <= 5)
);

-- Indexes
CREATE INDEX idx_courses_instructor ON courses(instructor_id);
CREATE INDEX idx_courses_category ON courses(category);
CREATE INDEX idx_courses_level ON courses(level);
CREATE INDEX idx_courses_status ON courses(status);
CREATE INDEX idx_courses_start_date ON courses(start_date);
CREATE INDEX idx_courses_rating ON courses(avg_rating DESC);
CREATE INDEX idx_courses_tags ON courses USING gin(tags);
CREATE INDEX idx_courses_format ON courses(format);

-- Comments
COMMENT ON TABLE courses IS 'Educational courses and programs';
COMMENT ON COLUMN courses.format IS 'self-paced, cohort-based, or live';
COMMENT ON COLUMN courses.curriculum IS 'JSONB field with course structure';

COMMIT;
