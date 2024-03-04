CREATE TABLE IF NOT EXISTS goods
(
    id          SERIAL PRIMARY KEY,
    project_id  INT  NOT NULL REFERENCES projects (id),
    name        TEXT NOT NULL,
    description TEXT,
    priority    INT  NOT NULL DEFAULT 1,
    removed     BOOLEAN       DEFAULT FALSE,
    created_at  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_priority CHECK (priority > 0)
);
CREATE INDEX IF NOT EXISTS project_id_idx ON goods (project_id);
CREATE INDEX IF NOT EXISTS name_idx ON goods (name);

CREATE OR REPLACE FUNCTION set_priority_on_insert()
    RETURNS TRIGGER AS
$set_priority_on_insert$
BEGIN
    NEW.priority = (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods);
    RETURN NEW;
END;
$set_priority_on_insert$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER set_priority_trigger
    BEFORE INSERT
    ON goods
    FOR EACH ROW
EXECUTE FUNCTION set_priority_on_insert();
