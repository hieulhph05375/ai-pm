-- 000018_refactor_issues_to_categories.up.sql

-- Add new ID columns
ALTER TABLE issues 
ADD COLUMN type_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
ADD COLUMN priority_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
ADD COLUMN status_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;

-- Note: We will handle data migration in a separate seed script or manual SQL
-- for more complex logic. For now, we drop the old columns as they are no longer needed
-- and were checked against the specific hardcoded strings.

ALTER TABLE issues DROP COLUMN issue_type;
ALTER TABLE issues DROP COLUMN priority;
ALTER TABLE issues DROP COLUMN status;

-- Drop old indices and create new ones
DROP INDEX IF EXISTS idx_issues_status;
CREATE INDEX IF NOT EXISTS idx_issues_type_id ON issues(type_id);
CREATE INDEX IF NOT EXISTS idx_issues_priority_id ON issues(priority_id);
CREATE INDEX IF NOT EXISTS idx_issues_status_id ON issues(status_id);
