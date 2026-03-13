-- 000026_granular_project_permissions.down.sql
BEGIN;

-- 1. Remove the granular permissions
DELETE FROM project_role_permissions 
WHERE project_permission_id IN (
    SELECT id FROM project_permissions 
    WHERE name IN (
        'project:wbs:create', 'project:wbs:update', 'project:wbs:delete',
        'project:risk:create', 'project:risk:update', 'project:risk:delete',
        'project:issue:create', 'project:issue:update', 'project:issue:delete',
        'project:team:create', 'project:team:update', 'project:team:delete',
        'project:roles:create', 'project:roles:update', 'project:roles:delete'
    )
);

DELETE FROM project_permissions 
WHERE name IN (
    'project:wbs:create', 'project:wbs:update', 'project:wbs:delete',
    'project:risk:create', 'project:risk:update', 'project:risk:delete',
    'project:issue:create', 'project:issue:update', 'project:issue:delete',
    'project:team:create', 'project:team:update', 'project:team:delete',
    'project:roles:create', 'project:roles:update', 'project:roles:delete'
);

-- 2. Restore older permissions
INSERT INTO project_permissions (name, description, module) VALUES
('project:team:manage', 'Add, remove, or change roles of project members', 'Team'),
('project:roles:manage', 'Create or update custom project roles', 'Team'),
('project:wbs:edit', 'Create, update, or delete WBS nodes', 'WBS'),
('project:risk:manage', 'Manage project risks', 'Risks'),
('project:issue:manage', 'Manage project issues', 'Issues')
ON CONFLICT (name) DO NOTHING;

COMMIT;
