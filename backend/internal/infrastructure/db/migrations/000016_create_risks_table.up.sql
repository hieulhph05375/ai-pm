CREATE TABLE IF NOT EXISTS risks (
    id              SERIAL PRIMARY KEY,
    project_id      INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    probability     SMALLINT NOT NULL DEFAULT 1 CHECK (probability BETWEEN 1 AND 5),
    impact          SMALLINT NOT NULL DEFAULT 1 CHECK (impact BETWEEN 1 AND 5),
    status          VARCHAR(50) NOT NULL DEFAULT 'Open' CHECK (status IN ('Open', 'Mitigated', 'Closed')),
    owner_id        INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_risks_project_id ON risks(project_id);
CREATE INDEX IF NOT EXISTS idx_risks_status ON risks(status);
