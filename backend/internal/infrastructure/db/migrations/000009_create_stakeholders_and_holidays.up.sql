-- Table: stakeholders
-- Stores organization-wide stakeholders (internal or external)
CREATE TABLE IF NOT EXISTS stakeholders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    organization VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Table: project_stakeholders
-- Maps stakeholders to projects with specific roles
CREATE TABLE IF NOT EXISTS project_stakeholders (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    stakeholder_id INTEGER NOT NULL REFERENCES stakeholders(id) ON DELETE CASCADE,
    project_role VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, stakeholder_id)
);

-- Table: holidays
-- Global holiday management (State & Company types)
CREATE TABLE IF NOT EXISTS holidays (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL, -- 'state' or 'company'
    is_recurring BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indices for performance
CREATE INDEX IF NOT EXISTS idx_holidays_date ON holidays(date);
CREATE INDEX IF NOT EXISTS idx_project_stakeholders_project ON project_stakeholders(project_id);
