CREATE TABLE IF NOT EXISTS projects
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO projects (name)
VALUES ('Первая запись')
