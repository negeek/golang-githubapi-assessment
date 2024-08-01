CREATE TABLE repositories (
    id SERIAL PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    language TEXT,
    forks_count INT,
    stars_count INT,
    open_issues_count INT,
    watchers_count INT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_repositories_name ON repositories(name);
