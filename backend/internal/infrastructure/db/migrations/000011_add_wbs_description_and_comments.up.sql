-- Add description field to wbs_nodes
ALTER TABLE wbs_nodes ADD COLUMN description TEXT;

-- Create wbs_comments table
CREATE TABLE wbs_comments (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    node_id INTEGER NOT NULL REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wbs_comments_node_id ON wbs_comments(node_id);
