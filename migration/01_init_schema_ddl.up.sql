-- Sequence Definition

DROP SEQUENCE IF EXISTS permission_seq;

CREATE SEQUENCE permission_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
	CACHE 1
	NO CYCLE;


-- role_seq definition

DROP SEQUENCE IF EXISTS role_seq;

CREATE SEQUENCE role_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
	CACHE 1
	NO CYCLE;


-- user_seq definition

DROP SEQUENCE IF EXISTS user_seq;

CREATE SEQUENCE user_seq
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    START 1
    CACHE 1
    NO CYCLE;

-- "permission" definition

-- Drop table

DROP TABLE IF EXISTS "permission";

CREATE TABLE "permission" (
    id int8 NOT NULL DEFAULT nextval('permission_seq'::regclass),
    "name" varchar(255) NULL,
    description varchar(255) NULL,
    created_at timestamp(6) NULL,
    updated_at timestamp(6) NULL,
    deleted_at timestamp(6) NULL,
    created_by int8 NULL,
    updated_by int8 NULL,
    deleted_by int8 NULL,
    CONSTRAINT permission_pkey PRIMARY KEY (id)
);

-- "role" definition

-- Drop table

DROP TABLE IF EXISTS "role";

CREATE TABLE "role" (
    id int8 NOT NULL DEFAULT nextval('role_seq'::regclass),
    "name" varchar(255) NULL,
    description varchar(255) NULL,
    created_at timestamp(6) NULL,
    updated_at timestamp(6) NULL,
    deleted_at timestamp(6) NULL,
    created_by int8 NULL,
    updated_by int8 NULL,
    deleted_by int8 NULL,
    CONSTRAINT role_pkey PRIMARY KEY (id)
);

-- "user" definition

-- Drop table

DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    id int8 NOT NULL DEFAULT nextval('user_seq'::regclass),
    username varchar(255) NULL,
    email varchar(255) NULL,
    "password" varchar(255) NULL,
    is_supper bool default false,
    last_ip varchar(255) NULL,
    last_login timestamp(6) NULL,
    created_at timestamp(6) NULL,
    updated_at timestamp(6) NULL,
    deleted_at timestamp(6) NULL,
    created_by int8 NULL,
    updated_by int8 NULL,
    deleted_by int8 NULL,
    is_deleted bool default false,
    is_active bool default true,
    is_locked bool default false,  -- Thêm cột is_locked
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

-- role_permission definition

-- Drop table

DROP TABLE IF EXISTS role_permission;

CREATE TABLE role_permission (
    permission_id int8 NOT NULL,
    role_id int8 NOT NULL,
    CONSTRAINT role_permission_pkey PRIMARY KEY (permission_id, role_id),
    CONSTRAINT "fk_role_permission_permission_id" FOREIGN KEY (permission_id) REFERENCES "permission"(id),
    CONSTRAINT "fk_role_permission_role_id" FOREIGN KEY (role_id) REFERENCES "role"(id)
);

-- user_role definition

-- Drop table

DROP TABLE IF EXISTS user_role;

CREATE TABLE user_role (
    user_id int8 NOT NULL,
    role_id int8 NOT NULL,
    CONSTRAINT user_role_pkey PRIMARY KEY (role_id, user_id),
    CONSTRAINT "fk_user_role_user_id" FOREIGN KEY (user_id) REFERENCES "user"(id),
    CONSTRAINT "fk_user_role_role_id" FOREIGN KEY (role_id) REFERENCES "role"(id)
);