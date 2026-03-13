-- 000025_migrate_project_member_roles.down.sql

ALTER TABLE project_members ADD COLUMN IF NOT EXISTS project_role VARCHAR(50);

UPDATE project_members pm
SET project_role = pr.name
FROM project_roles pr
WHERE pm.project_role_id = pr.id;

ALTER TABLE project_members ALTER COLUMN project_role_id DROP NOT NULL;
