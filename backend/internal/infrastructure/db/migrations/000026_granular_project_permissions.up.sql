-- 000026_granular_project_permissions.up.sql
BEGIN;

-- 1. Ensure old grouped permissions are gone (in case they were re-inserted by a different script)
DELETE FROM project_role_permissions 
WHERE project_permission_id IN (
    SELECT id FROM project_permissions 
    WHERE name IN ('project:wbs:edit', 'project:risk:manage', 'project:issue:manage', 'project:team:manage', 'project:roles:manage')
);

DELETE FROM project_permissions 
WHERE name IN ('project:wbs:edit', 'project:risk:manage', 'project:issue:manage', 'project:team:manage', 'project:roles:manage');

-- 2. Insert or Update granular permissions
INSERT INTO project_permissions (name, description, module) VALUES
('project:wbs:create', 'Create new WBS nodes', 'WBS'),
('project:wbs:update', 'Update existing WBS nodes', 'WBS'),
('project:wbs:delete', 'Delete WBS nodes', 'WBS'),
('project:risk:create', 'Log new project risks', 'Risks'),
('project:risk:update', 'Update existing risks', 'Risks'),
('project:risk:delete', 'Delete existing risks', 'Risks'),
('project:issue:create', 'Log new issues or bugs', 'Issues'),
('project:issue:update', 'Update existing issues', 'Issues'),
('project:issue:delete', 'Delete issues', 'Issues'),
('project:team:create', 'Add members to the project', 'Team'),
('project:team:update', 'Change member roles', 'Team'),
('project:team:delete', 'Remove members from the project', 'Team'),
('project:roles:create', 'Create custom project roles', 'Roles & Permissions'),
('project:roles:update', 'Update project roles and permissions', 'Roles & Permissions'),
('project:roles:delete', 'Delete project roles', 'Roles & Permissions')
ON CONFLICT (name) DO UPDATE SET 
    description = EXCLUDED.description,
    module = EXCLUDED.module;

-- 2. Migrate existing mappings for WBS
INSERT INTO project_role_permissions (project_role_id, project_permission_id)
SELECT prp.project_role_id, p_new.id
FROM project_role_permissions prp
JOIN project_permissions p_old ON prp.project_permission_id = p_old.id
JOIN project_permissions p_new ON p_new.name IN ('project:wbs:create', 'project:wbs:update', 'project:wbs:delete')
WHERE p_old.name = 'project:wbs:edit'
ON CONFLICT DO NOTHING;

-- 3. Migrate existing mappings for Risks
INSERT INTO project_role_permissions (project_role_id, project_permission_id)
SELECT prp.project_role_id, p_new.id
FROM project_role_permissions prp
JOIN project_permissions p_old ON prp.project_permission_id = p_old.id
JOIN project_permissions p_new ON p_new.name IN ('project:risk:create', 'project:risk:update', 'project:risk:delete')
WHERE p_old.name = 'project:risk:manage'
ON CONFLICT DO NOTHING;

-- 4. Migrate existing mappings for Issues
INSERT INTO project_role_permissions (project_role_id, project_permission_id)
SELECT prp.project_role_id, p_new.id
FROM project_role_permissions prp
JOIN project_permissions p_old ON prp.project_permission_id = p_old.id
JOIN project_permissions p_new ON p_new.name IN ('project:issue:create', 'project:issue:update', 'project:issue:delete')
WHERE p_old.name = 'project:issue:manage'
ON CONFLICT DO NOTHING;

-- 5. Migrate existing mappings for Team
INSERT INTO project_role_permissions (project_role_id, project_permission_id)
SELECT prp.project_role_id, p_new.id
FROM project_role_permissions prp
JOIN project_permissions p_old ON prp.project_permission_id = p_old.id
JOIN project_permissions p_new ON p_new.name IN ('project:team:create', 'project:team:update', 'project:team:delete')
WHERE p_old.name = 'project:team:manage'
ON CONFLICT DO NOTHING;

-- 6. Migrate existing mappings for Roles
INSERT INTO project_role_permissions (project_role_id, project_permission_id)
SELECT prp.project_role_id, p_new.id
FROM project_role_permissions prp
JOIN project_permissions p_old ON prp.project_permission_id = p_old.id
JOIN project_permissions p_new ON p_new.name IN ('project:roles:create', 'project:roles:update', 'project:roles:delete')
WHERE p_old.name = 'project:roles:manage'
ON CONFLICT DO NOTHING;

COMMIT;

