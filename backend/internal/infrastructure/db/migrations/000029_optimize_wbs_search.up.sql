-- Enable pg_trgm extension for trigram search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create GIN trigram index on title for fast ILIKE searches
CREATE INDEX IF NOT EXISTS idx_wbs_nodes_title_trgm ON wbs_nodes USING gin (title gin_trgm_ops);
