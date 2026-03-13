-- 000031_add_role_id_to_project_stakeholders.up.sql
ALTER TABLE project_stakeholders ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES categories(id);

-- Update existing data: try to map project_role string to category name
-- (Best effort mapping)
UPDATE project_stakeholders ps 
SET role_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'STAKEHOLDER_ROLE') AND LOWER(name) = LOWER(ps.project_role))
WHERE ps.project_role IS NOT NULL;
