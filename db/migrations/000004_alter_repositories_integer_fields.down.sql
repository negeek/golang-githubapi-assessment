ALTER TABLE repositories
ALTER COLUMN forks_count TYPE INT USING forks_count::INT,
ALTER COLUMN stars_count TYPE INT USING stars_count::INT,
ALTER COLUMN open_issues_count TYPE INT USING open_issues_count::INT,
ALTER COLUMN watchers_count TYPE INT USING watchers_count::INT;