INSERT INTO "user"
(avatar, created_at, created_by, is_supper, last_ip, last_login, "password", is_active, updated_at, updated_by, username, email)
VALUES('', '2024-08-13 13:12:18.098', 0, true, '', '2024-08-13 13:12:18.098', '$2a$10$iLJXuE6KvxNqh193OpmYbuB4UJg3V4HAvulk5JxHpBffQgVO7PY3m', true, '2024-08-13 13:12:18.098', 0, 'admin', 'truongvu.pham@gmail.vn');
INSERT INTO "permission"
(created_at, created_by, description, "name", updated_at, updated_by)
VALUES('2024-08-13 13:12:18.098', 1, 'Create user', 'CREATE_USER', '2024-08-13 13:12:18.098', 1);
INSERT INTO "role"
(created_at, created_by, description, "name", updated_at, updated_by)
VALUES('2024-08-13 13:12:18.098', 1, 'Admin User', 'Administrator', '2024-08-13 13:12:18.098', 1);
INSERT INTO "role"
(created_at, created_by, description, "name", updated_at, updated_by)
VALUES('2024-08-13 13:12:18.098', 1, 'Default User', 'DEFAULT_USER', '2024-08-13 13:12:18.098', 1);
INSERT INTO role_permission
(permission_id, role_id)
VALUES(1, 1);