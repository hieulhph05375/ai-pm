-- 000007_add_is_admin_to_users.up.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE;
UPDATE users SET is_admin = TRUE WHERE role_id = 1; -- Assume Role ID 1 is Admin
