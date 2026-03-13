ALTER TABLE wbs_nodes ADD COLUMN estimated_effort NUMERIC(10, 2) DEFAULT 0;
ALTER TABLE wbs_nodes ADD COLUMN actual_effort NUMERIC(10, 2) DEFAULT 0;
