-- 000024_create_project_rbac_tables.down.sql

DROP TRIGGER IF EXISTS tr_project_roles_updated_at ON project_roles;
DROP FUNCTION IF EXISTS update_project_roles_updated_at();

ALTER TABLE project_members DROP COLUMN IF EXISTS project_role_id;
DROP TABLE IF EXISTS project_role_permissions;
DROP TABLE IF EXISTS project_roles;
DROP TABLE IF EXISTS project_permissions;
