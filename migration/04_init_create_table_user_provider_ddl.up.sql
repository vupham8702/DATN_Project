CREATE SEQUENCE user_prv_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE user_provider (
    id int8 NOT NULL DEFAULT nextval('user_prv_seq'::regclass),
    email varchar(255) NULL,
    provider varchar(255) NOT NULL,
    provider_identify varchar(255) NOT NULL,
    user_type varchar(255) NULL,
    account_info JSONB,
    user_id int8 NOT NULL UNIQUE,
    created_at timestamp(6) NULL,
    updated_at timestamp(6) NULL,
    deleted_at timestamp(6) NULL,
    created_by int8 NULL,
    updated_by int8 NULL,
    deleted_by int8 NULL,
    is_deleted bool default false,
    is_approved bool default false,
    approved_by int8 NULL,
    approval_note varchar(255) NULL,
    received_noti bool default true,
    CONSTRAINT user_prv_pkey PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
