-- 000018_refactor_issues_to_categories.down.sql

ALTER TABLE issues 
ADD COLUMN issue_type VARCHAR(50) NOT NULL DEFAULT 'Issue' CHECK (issue_type IN ('Bug', 'Issue')),
ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'Open' CHECK (status IN ('Open', 'In Progress', 'Resolved', 'Closed')),
ADD COLUMN priority VARCHAR(50) NOT NULL DEFAULT 'Medium' CHECK (priority IN ('Low', 'Medium', 'High', 'Critical'));

DROP INDEX IF EXISTS idx_issues_type_id;
DROP INDEX IF EXISTS idx_issues_priority_id;
DROP INDEX IF EXISTS idx_issues_status_id;

ALTER TABLE issues DROP COLUMN type_id;
ALTER TABLE issues DROP COLUMN priority_id;
ALTER TABLE issues DROP COLUMN status_id;

CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
