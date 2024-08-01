CREATE TABLE commits (
    id SERIAL PRIMARY KEY,
    sha TEXT UNIQUE NOT NULL,
    repo VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    message TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (repo) REFERENCES repositories(name)
);

CREATE INDEX idx_commits_sha ON commits(sha);
CREATE INDEX idx_commits_date ON commits(date);
CREATE INDEX idx_commits_repo ON commits(repo);

