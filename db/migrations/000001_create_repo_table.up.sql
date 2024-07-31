CREATE TABLE repositories (
    id SERIAL PRIMARY KEY,
    owner VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    description TEXT,
    url VARCHAR NOT NULL,
    language VARCHAR,
    forks_count INT,
    stars_count INT,
    open_issues_count INT,
    watchers_count INT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_repositories_name ON repositories(name);
