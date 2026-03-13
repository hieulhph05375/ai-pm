-- ==============================================================================
-- Project Management System - Production Seed Data
-- Run AFTER migrations have been applied.
-- All statements are idempotent (safe to re-run).
-- ==============================================================================

-- 1. Default Administrator
-- Password: password (CHANGE THIS IN PRODUCTION)
INSERT INTO users (email, hashed_password, full_name, role_id, is_active, is_admin)
SELECT 'admin@admin.com', '$2a$10$nFVpvMZgQm1Dd/RQaDzXCe5XbH2QX82oe0Iy8mgiXOwvbn4JaoTV.', 'System Administrator', 1, true, true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@admin.com');

-- ==============================================================================
-- 2. Category Types
-- ==============================================================================
INSERT INTO category_types (name, code, description) VALUES
-- Issue Management
('Issue Type',          'issue_type',          'Types of project issues'),
('Issue Priority',      'issue_priority',       'Priority levels for issues'),
('Issue Status',        'issue_status',         'Status workflow for issues'),
-- Project Management
('Project Status',      'PROJECT_STATUS',       'Standard statuses for project lifecycle'),
('Project Phase',       'PROJECT_PHASE',        'Standard phases for project management'),
('Portfolio Category',  'PORTFOLIO_CATEGORY',   'Strategic classification for projects'),
('Project Health',      'PROJECT_HEALTH',       'Visual health indicators (RAG)'),
('Priority Level',      'PRIORITY_LEVEL',       'Priority levels for projects'),
-- Task Management
('Task Status',         'TASK_STATUS',          'Simple workflow for personal tasks'),
('Task Priority',       'TASK_PRIORITY',        'Priority levels for personal tasks'),
-- Risk Management
('Risk Status',         'RISK_STATUS',          'Workflow for project risks'),
-- Stakeholders
('Stakeholder Role',    'STAKEHOLDER_ROLE',     'Organizational roles for stakeholders'),
-- Holidays
('Holiday Type',        'HOLIDAY_TYPE',         'Classification for holidays'),
-- WBS
('WBS Node Type',       'WBS_NODE_TYPE',        'Classification for WBS elements (Phase, Milestone, Task)')
ON CONFLICT (code) DO NOTHING;

-- ==============================================================================
-- 3. Categories
-- ==============================================================================

-- Issue Types
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Bug', '#ef4444'), ('Task', '#3b82f6'), ('Improvement', '#10b981')) AS v(name, color)
WHERE ct.code = 'issue_type'
ON CONFLICT (type_id, name) DO NOTHING;

-- Issue Priorities
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Low', '#94a3b8'), ('Medium', '#f59e0b'), ('High', '#ef4444'), ('Critical', '#7f1d1d')) AS v(name, color)
WHERE ct.code = 'issue_priority'
ON CONFLICT (type_id, name) DO NOTHING;

-- Issue Statuses
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Open', '#64748b'), ('In Progress', '#3b82f6'), ('Resolved', '#10b981'), ('Closed', '#1e293b')) AS v(name, color)
WHERE ct.code = 'issue_status'
ON CONFLICT (type_id, name) DO NOTHING;

-- Project Statuses
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Planned', 'slate'), ('Active', 'blue'), ('On Hold', 'amber'), ('Completed', 'emerald'), ('Cancelled', 'rose')) AS v(name, color)
WHERE ct.code = 'PROJECT_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Project Phases
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Initiation', 'slate'), ('Planning', 'blue'), ('Execution', 'violet'), ('Monitoring', 'amber'), ('Closing', 'emerald')) AS v(name, color)
WHERE ct.code = 'PROJECT_PHASE'
ON CONFLICT (type_id, name) DO NOTHING;

-- Portfolio Categories
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Strategic', 'violet'), ('Operational', 'blue'), ('Compliance', 'amber')) AS v(name, color)
WHERE ct.code = 'PORTFOLIO_CATEGORY'
ON CONFLICT (type_id, name) DO NOTHING;

-- Project Health (RAG)
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Green', 'emerald'), ('Yellow', 'amber'), ('Red', 'rose')) AS v(name, color)
WHERE ct.code = 'PROJECT_HEALTH'
ON CONFLICT (type_id, name) DO NOTHING;

-- Task Statuses
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Todo', 'slate'), ('In Progress', 'blue'), ('Done', 'emerald')) AS v(name, color)
WHERE ct.code = 'TASK_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Task Priorities
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Low', 'slate'), ('Medium', 'blue'), ('High', 'amber'), ('Urgent', 'rose')) AS v(name, color)
WHERE ct.code = 'TASK_PRIORITY'
ON CONFLICT (type_id, name) DO NOTHING;

-- Risk Statuses
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Open', 'rose'), ('Mitigated', 'blue'), ('Closed', 'emerald')) AS v(name, color)
WHERE ct.code = 'RISK_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Stakeholder Roles
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Internal', 'blue'), ('External', 'violet'), ('Vendor', 'amber'), ('Government', 'slate')) AS v(name, color)
WHERE ct.code = 'STAKEHOLDER_ROLE'
ON CONFLICT (type_id, name) DO NOTHING;

-- Holiday Types
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('State', 'rose'), ('Company', 'blue')) AS v(name, color)
WHERE ct.code = 'HOLIDAY_TYPE'
ON CONFLICT (type_id, name) DO NOTHING;

-- WBS Node Types
INSERT INTO categories (type_id, name, color)
SELECT ct.id, v.name, v.color FROM category_types ct
CROSS JOIN (VALUES ('Phase', 'slate'), ('Milestone', 'amber'), ('Work Package', 'blue')) AS v(name, color)
WHERE ct.code = 'WBS_NODE_TYPE'
ON CONFLICT (type_id, name) DO NOTHING;

-- ==============================================================================
-- 4. Vietnam Public Holidays 2026 (Official)
-- ==============================================================================
INSERT INTO holidays (name, date, type, is_recurring) VALUES
('New Year Day',                   '2026-01-01', 'state', true),
('Lunar New Year Eve (Tết)',        '2026-02-16', 'state', false),
('Lunar New Year Day 1',           '2026-02-17', 'state', false),
('Lunar New Year Day 2',           '2026-02-18', 'state', false),
('Lunar New Year Day 3',           '2026-02-19', 'state', false),
('Lunar New Year Day 4',           '2026-02-20', 'state', false),
('Lunar New Year Day 5',           '2026-02-21', 'state', false),
('Lunar New Year Day 6',           '2026-02-22', 'state', false),
('Hung Kings Commemoration Day',   '2026-04-26', 'state', false),
('Reunification Day',              '2026-04-30', 'state', true),
('International Labor Day',        '2026-05-01', 'state', true),
('Vietnam National Day',           '2026-09-02', 'state', true),
('Vietnam National Day (Bridge)',   '2026-09-03', 'state', false)
ON CONFLICT (date) DO NOTHING;

-- ==============================================================================
-- 5. System Settings
-- ==============================================================================
INSERT INTO system_settings (key, value) VALUES
('org_name',           '"Project Management System"'::jsonb),
('default_currency',   '"VND"'::jsonb),
('timezone',           '"UTC+7"'::jsonb),
('work_hours_start',   '"08:00"'::jsonb),
('work_hours_end',     '"17:00"'::jsonb)
ON CONFLICT (key) DO NOTHING;
