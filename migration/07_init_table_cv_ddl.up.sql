CREATE TABLE IF NOT EXISTS datn_backend.cv_template (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    thumbnail_path VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS datn_backend.user_cv (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES datn_backend."user"(id),
    template_id INTEGER REFERENCES datn_backend.cv_template(id),
    file_path VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_by INTEGER,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_by INTEGER,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

    );

-- Tạo index để tăng tốc truy vấn
CREATE INDEX IF NOT EXISTS idx_user_cv_user_id ON datn_backend.user_cv(user_id);
CREATE INDEX IF NOT EXISTS idx_user_cv_template_id ON datn_backend.user_cv(template_id);