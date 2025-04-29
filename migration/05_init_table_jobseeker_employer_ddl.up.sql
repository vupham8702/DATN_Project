CREATE SEQUENCE IF NOT EXISTS jobseeker_profile_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE IF NOT EXISTS jobseeker_profile (
    id int8 NOT NULL DEFAULT nextval('jobseeker_profile_seq'::regclass),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    created_by int8 NOT NULL DEFAULT 0,
    updated_by int8 NOT NULL DEFAULT 0,
    deleted_by int8 NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    user_id int8 NOT NULL UNIQUE,
    date_of_birth TIMESTAMP,
    gender VARCHAR(20),
    phone_number VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    profile_title VARCHAR(255),
    about TEXT,
    skills TEXT,
    education JSONB,
    experience JSONB,
    certifications JSONB,
    languages JSONB,
    resume_url TEXT,
    profile_picture TEXT,
    profile_complete BOOLEAN NOT NULL DEFAULT FALSE,
    availability VARCHAR(50),
    linkedin_profile VARCHAR(255),
    github_profile VARCHAR(255),
    website_url VARCHAR(255),
    CONSTRAINT jobseeker_profile_pkey PRIMARY KEY (id),
    CONSTRAINT fk_jobseeker_profile_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);


CREATE SEQUENCE employer_profile_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE IF NOT EXISTS employer_profile (
    id int8 NOT NULL DEFAULT nextval('employer_profile_seq'::regclass),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    created_by int8 NOT NULL DEFAULT 0,
    updated_by int8 NOT NULL DEFAULT 0,
    deleted_by int8 NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    user_id int8 NOT NULL UNIQUE,
    company_name VARCHAR(255),
    company_size VARCHAR(50),
    industry VARCHAR(100),
    company_logo TEXT,
    company_banner TEXT,
    website VARCHAR(255),
    founded VARCHAR(255),
    about TEXT,
    mission TEXT,
    phone_number VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100),
    facebook_url VARCHAR(255),
    twitter_url VARCHAR(255),
    linkedin_url VARCHAR(255),
    benefits JSONB,
    culture JSONB,
    profile_complete BOOLEAN NOT NULL DEFAULT FALSE,
    tax_code VARCHAR(50),
    business_license VARCHAR(100),
    contact_person_name VARCHAR(255),
    contact_person_role VARCHAR(100),

    CONSTRAINT employer_profile_pkey PRIMARY KEY (id),
    CONSTRAINT fk_employer_profile_user FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

