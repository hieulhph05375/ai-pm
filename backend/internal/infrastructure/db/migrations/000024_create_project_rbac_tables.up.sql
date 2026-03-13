-- 000024_create_project_rbac_tables.up.sql

-- 1. Project Permissions (System-defined available actions for projects)
CREATE TABLE IF NOT EXISTS project_permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL, -- e.g., 'project:wbs:edit'
    description TEXT,
    module VARCHAR(50) NOT NULL,       -- e.g., 'WBS', 'Risks', 'Team'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Project Roles (Defined per project)
CREATE TABLE IF NOT EXISTS project_roles (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color VARCHAR(20) DEFAULT '#64748B',
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, name)
);

-- 3. Project Role Permissions (Mapping roles to permissions)
CREATE TABLE IF NOT EXISTS project_role_permissions (
    project_role_id INTEGER REFERENCES project_roles(id) ON DELETE CASCADE,
    project_permission_id INTEGER REFERENCES project_permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(project_role_id, project_permission_id)
);

-- 4. Initial Project Permissions
INSERT INTO project_permissions (name, description, module) VALUES
('project:read', 'View basic project information', 'Project'),
('project:update', 'Update project configuration and metadata', 'Project'),
('project:team:view', 'View project team members', 'Team'),
('project:team:manage', 'Add, remove, or change roles of project members', 'Team'),
('project:roles:manage', 'Create or update custom project roles', 'Team'),
('project:wbs:view', 'View project WBS and Gantt chart', 'WBS'),
('project:wbs:edit', 'Create, update, or delete WBS nodes', 'WBS'),
('project:risk:view', 'View project risk register', 'Risks'),
('project:risk:manage', 'Manage project risks', 'Risks'),
('project:issue:view', 'View project issue list', 'Issues'),
('project:issue:manage', 'Manage project issues', 'Issues'),
('project:timesheet:view', 'View timesheets for this project', 'Timesheets'),
('project:timesheet:approve', 'Approve or reject project timesheets', 'Timesheets'),
('project:pmi:view', 'View EVM and performance metrics', 'Reports')
ON CONFLICT (name) DO NOTHING;

-- 5. Alter project_members to use project_role_id
ALTER TABLE project_members ADD COLUMN IF NOT EXISTS project_role_id INTEGER REFERENCES project_roles(id);

-- 6. Trigger for updated_at in project_roles
CREATE OR REPLACE FUNCTION update_project_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS tr_project_roles_updated_at ON project_roles;
CREATE TRIGGER tr_project_roles_updated_at
    BEFORE UPDATE ON project_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_project_roles_updated_at();
