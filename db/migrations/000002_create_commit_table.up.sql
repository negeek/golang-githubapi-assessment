CREATE TABLE commits (
    id SERIAL PRIMARY KEY,
    sha VARCHAR NOT NULL,
    repo INTEGER NOT NULL,
    author_name VARCHAR NOT NULL,
    author_email VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    message TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (repo) REFERENCES repositories(id)
);

CREATE INDEX idx_commits_sha ON commits(sha);
