-- Tạo sequence cho bảng post_job
CREATE SEQUENCE IF NOT EXISTS post_job_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Tạo bảng post_job
CREATE TABLE IF NOT EXISTS post_job (
    id int8 NOT NULL DEFAULT nextval('post_job_seq'::regclass),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    created_by int8 NOT NULL DEFAULT 0,
    updated_by int8 NOT NULL DEFAULT 0,
    deleted_by int8 NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Thông tin bài đăng
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    logo TEXT,
    location VARCHAR(255),
    salary VARCHAR(100),
    status VARCHAR(50) DEFAULT 'pending',
    type VARCHAR(50),
    time_frame VARCHAR(100),
    experience VARCHAR(100),
    gender VARCHAR(20),
    description TEXT NOT NULL,
    applications_count INT DEFAULT 0,
    
    -- Các trường bổ sung
    requirements TEXT,
    benefits TEXT,
    deadline TIMESTAMP,
    positions INT DEFAULT 1,
    views INT DEFAULT 0,
    
    CONSTRAINT post_job_pkey PRIMARY KEY (id)
);

-- Tạo index để tìm kiếm nhanh
CREATE INDEX IF NOT EXISTS idx_post_job_created_by ON post_job(created_by);
CREATE INDEX IF NOT EXISTS idx_post_job_status ON post_job(status);
CREATE INDEX IF NOT EXISTS idx_post_job_title ON post_job(title);
CREATE INDEX IF NOT EXISTS idx_post_job_location ON post_job(location);
CREATE INDEX IF NOT EXISTS idx_post_job_type ON post_job(type);
CREATE INDEX IF NOT EXISTS idx_post_job_is_deleted ON post_job(is_deleted);

-- Tạo bảng post_job_application để lưu thông tin ứng tuyển
CREATE SEQUENCE IF NOT EXISTS post_job_application_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS post_job_application (
    id int8 NOT NULL DEFAULT nextval('post_job_application_seq'::regclass),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    created_by int8 NOT NULL DEFAULT 0, -- user_id của người ứng tuyển
    updated_by int8 NOT NULL DEFAULT 0,
    deleted_by int8 NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    
    post_job_id int8 NOT NULL,
    resume_url TEXT NOT NULL,
    cover_letter TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- pending, reviewed, shortlisted, interviewed, offered, rejected
    notes TEXT,
    
    CONSTRAINT post_job_application_pkey PRIMARY KEY (id),
    CONSTRAINT fk_post_job_application_post_job FOREIGN KEY (post_job_id) REFERENCES post_job(id) ON DELETE CASCADE
);

-- Tạo index cho bảng application
CREATE INDEX IF NOT EXISTS idx_post_job_application_post_job_id ON post_job_application(post_job_id);
CREATE INDEX IF NOT EXISTS idx_post_job_application_created_by ON post_job_application(created_by);
CREATE INDEX IF NOT EXISTS idx_post_job_application_status ON post_job_application(status);