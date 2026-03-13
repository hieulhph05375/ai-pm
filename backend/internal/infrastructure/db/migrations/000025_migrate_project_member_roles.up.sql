-- 000025_migrate_project_member_roles.up.sql

DO $$
DECLARE
    p_id INTEGER;
    role_id INTEGER;
    perm_id INTEGER;
BEGIN
    -- 1. Create Default Roles for ALL existing projects
    FOR p_id IN SELECT id FROM projects LOOP
        
        -- Lead Role
        INSERT INTO project_roles (project_id, name, description, color, is_default)
        VALUES (p_id, 'Lead', 'Full project management permissions', '#EF4444', true)
        ON CONFLICT (project_id, name) DO NOTHING
        RETURNING id INTO role_id;

        IF role_id IS NOT NULL THEN
            -- Assign ALL permissions to Lead
            INSERT INTO project_role_permissions (project_role_id, project_permission_id)
            SELECT role_id, id FROM project_permissions ON CONFLICT DO NOTHING;
        END IF;

        -- Member Role
        INSERT INTO project_roles (project_id, name, description, color, is_default)
        VALUES (p_id, 'Member', 'Standard project contributor', '#3B82F6', false)
        ON CONFLICT (project_id, name) DO NOTHING
        RETURNING id INTO role_id;

        IF role_id IS NOT NULL THEN
            -- Assign subset of permissions to Member
            INSERT INTO project_role_permissions (project_role_id, project_permission_id)
            SELECT role_id, id FROM project_permissions 
            WHERE name IN ('project:read', 'project:team:view', 'project:wbs:view', 'project:risk:view', 'project:issue:view', 'project:timesheet:view')
            ON CONFLICT DO NOTHING;
        END IF;

        -- Observer Role
        INSERT INTO project_roles (project_id, name, description, color, is_default)
        VALUES (p_id, 'Observer', 'Read-only access to project data', '#64748B', false)
        ON CONFLICT (project_id, name) DO NOTHING
        RETURNING id INTO role_id;

        IF role_id IS NOT NULL THEN
            -- Assign read-only permissions to Observer
            INSERT INTO project_role_permissions (project_role_id, project_permission_id)
            SELECT role_id, id FROM project_permissions 
            WHERE name IN ('project:read', 'project:team:view', 'project:wbs:view', 'project:risk:view', 'project:issue:view')
            ON CONFLICT DO NOTHING;
        END IF;

    END LOOP;

    -- 2. Map existing project_members to the new roles
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'project_members' AND column_name = 'project_role') THEN
        UPDATE project_members pm
        SET project_role_id = pr.id
        FROM project_roles pr
        WHERE pm.project_id = pr.project_id 
        AND (
            (pm.project_role = 'Lead' AND pr.name = 'Lead') OR
            (pm.project_role = 'Member' AND pr.name = 'Member') OR
            (pm.project_role = 'Observer' AND pr.name = 'Observer') OR
            (pm.project_role NOT IN ('Lead', 'Member', 'Observer') AND pr.is_default = true)
        );

        -- 3. Cleanup: assign default 'Member' role
        UPDATE project_members pm
        SET project_role_id = pr.id
        FROM project_roles pr
        WHERE pm.project_role_id IS NULL 
        AND pm.project_id = pr.project_id
        AND pr.name = 'Member';
    END IF;

END $$;

-- 4. Final Refactoring
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'project_members' AND column_name = 'project_role_id' AND is_nullable = 'YES') THEN
        ALTER TABLE project_members ALTER COLUMN project_role_id SET NOT NULL;
    END IF;
END $$;
ALTER TABLE project_members DROP COLUMN IF EXISTS project_role;
