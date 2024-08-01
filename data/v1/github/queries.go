package github

import (
	"context"
	"fmt"
	"time"

	"github.com/negeek/golang-githubapi-assessment/db"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func (c *Commit) FindLatestRepoCommitByDate() error {
	query := `
		SELECT id, sha, repo, author_name, author_email, url, message, date
		FROM commits
		WHERE repo = $1
		ORDER BY date DESC
		LIMIT 1;
	`

	row := db.PostgreSQLDB.QueryRow(context.Background(), query, c.Repo)

	err := row.Scan(&c.ID, &c.SHA, &c.Repo, &c.AuthorName, &c.AuthorEmail, &c.URL, &c.Message, &c.Date)
	if err != nil {
		return err
	}

	return nil
}

func CreateCommits(commits []Commit) error {
	for _, c := range commits {
		utils.Time(c, true)
		query := "INSERT INTO commits (id, sha, repo, author_name, author_email, url, message, date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, c.ID, c.SHA, c.Repo, c.AuthorName, c.AuthorEmail, c.URL, c.Message, c.Date)
		if err != nil {
			return err
		}
	}
	return nil

}

func (r *Repository) Create() error {
	utils.Time(r, true)
	query := `
		INSERT INTO repositories (
			id, 
			owner, 
			name, 
			description, 
			url, 
			language, 
			forks_count, 
			stars_count, 
			open_issues_count, 
			watchers_count, 
			created_at, 
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`
	_, err := db.PostgreSQLDB.Exec(
		context.Background(),
		query,
		r.ID,
		r.Owner,
		r.Name,
		r.Description,
		r.URL,
		r.Language,
		r.ForksCount,
		r.StarsCount,
		r.OpenIssuesCount,
		r.WatchersCount,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SetupData) Create() error {
	utils.Time(s, true)
	query := "INSERT INTO setup_data (id, owner, repo, from_date, to_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.PostgreSQLDB.Exec(context.Background(), query, s.ID, s.Owner, s.Repo, s.FromDate, s.ToDate, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func Get_all_setup_data() ([]SetupData, error) {
	query := "SELECT id, owner, repo, from_date, to_date, created_at, updated_at FROM setup_data"
	rows, err := db.PostgreSQLDB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var setupData []SetupData
	for rows.Next() {
		var s SetupData
		err := rows.Scan(&s.ID, &s.Owner, &s.Repo, &s.FromDate, &s.ToDate, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		setupData = append(setupData, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return setupData, nil
}

func GetTopNCommitAuthors(repo string, topN int) ([]map[string]interface{}, error) {
	query := `
		SELECT author_name, COUNT(*) as commit_count
		FROM commits
		WHERE repo = $1
		GROUP BY author_name
		ORDER BY commit_count DESC
		LIMIT $2;
	`

	rows, err := db.PostgreSQLDB.Query(context.Background(), query, repo, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		authorName  string
		commitCount int
		authors     []map[string]interface{}
	)

	for rows.Next() {

		if err := rows.Scan(&authorName, &commitCount); err != nil {
			return nil, err
		}

		author := map[string]interface{}{
			"name":         authorName,
			"commit_count": commitCount,
		}

		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func GetCommitsByRepoName(repo string) ([]map[string]interface{}, error) {
	query := `
		SELECT sha, repo, author_name, author_email, url, message, date
		FROM commits
		WHERE repo = $1
		ORDER BY date DESC;
	`

	rows, err := db.PostgreSQLDB.Query(context.Background(), query, repo)
	if err != nil {
		return nil, fmt.Errorf("error querying commits for repository %s: %w", repo, err)
	}
	defer rows.Close()

	var commits []map[string]interface{}

	for rows.Next() {
		var (
			sha         string
			repo        string
			authorName  string
			authorEmail string
			url         string
			message     string
			date        time.Time
		)
		if err := rows.Scan(&sha, &repo, &authorName, &authorEmail, &url, &message, &date); err != nil {
			return nil, err
		}

		commit := map[string]interface{}{
			"sha":          sha,
			"repo":         repo,
			"author_name":  authorName,
			"author_email": authorEmail,
			"url":          url,
			"message":      message,
			"date":         date,
		}

		commits = append(commits, commit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commits, nil
}
