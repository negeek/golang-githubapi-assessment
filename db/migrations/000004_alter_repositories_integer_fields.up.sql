ALTER TABLE repositories
ALTER COLUMN forks_count TYPE FLOAT USING forks_count::FLOAT,
ALTER COLUMN stars_count TYPE FLOAT USING stars_count::FLOAT,
ALTER COLUMN open_issues_count TYPE FLOAT USING open_issues_count::FLOAT,
ALTER COLUMN watchers_count TYPE FLOAT USING watchers_count::FLOAT;