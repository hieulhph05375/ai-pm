-- 000005_remove_username.up.sql

ALTER TABLE users DROP COLUMN IF EXISTS username;
