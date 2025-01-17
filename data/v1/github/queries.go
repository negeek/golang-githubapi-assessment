package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/negeek/golang-githubapi-assessment/db"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func (c *Commit) FindLatestRepoCommitByDate() (bool, error) {
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
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CreateCommits(commits []Commit) {
	for _, c := range commits {
		utils.Time(&c, true)
		query := "INSERT INTO commits (sha, repo, author_name, author_email, url, message, date) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, c.SHA, strings.ToLower(c.Repo), c.AuthorName, c.AuthorEmail, c.URL, c.Message, c.Date)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (r *Repository) Create() error {
	utils.Time(r, true)
	query := `
		INSERT INTO repositories (
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
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`
	_, err := db.PostgreSQLDB.Exec(
		context.Background(),
		query,
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

func FindRepoByName(name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM repositories WHERE name=$1)"
	name = strings.ToLower(name)
	err := db.PostgreSQLDB.QueryRow(context.Background(), query, name).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (s *SetupData) CreateOrUpdate() error {
	utils.Time(s, true)
	query := `
		INSERT INTO setup_data (repo, from_date, to_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (repo) 
		DO UPDATE SET
			from_date = EXCLUDED.from_date,
			to_date = EXCLUDED.to_date,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
	`
	_, err := db.PostgreSQLDB.Exec(context.Background(), query,
		strings.ToLower(s.Repo),
		s.FromDate,
		s.ToDate,
		s.CreatedAt,
		s.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func SetEnvSetupData() {
	var err error

	s := &SetupData{
		Repo:     os.Getenv("GITHUB_REPO"),
		FromDate: time.Time{}, // default value
		ToDate:   time.Now(),  // default value
	}

	// Handle date fields
	s.FromDate, err = utils.HandleDateField(os.Getenv("GITHUB_COMMIT_FROM_DATE"), s.FromDate)
	if err != nil {
		log.Printf("invalid date format: %v", err)
		return
	}

	s.ToDate, err = utils.HandleDateField(os.Getenv("GITHUB_COMMIT_TO_DATE"), s.ToDate)
	if err != nil {
		log.Printf("invalid date format: %v", err)
		return
	}

	err = s.Validate()
	if err != nil {
		log.Println(err)
		return
	}

	err = s.CreateOrUpdate()
	if err != nil {
		log.Println(err)
	}
}

func GetAllSetUpData() ([]SetupData, error) {
	query := "SELECT id, repo, from_date, to_date, created_at, updated_at FROM setup_data"
	rows, err := db.PostgreSQLDB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var setupData []SetupData
	for rows.Next() {
		var s SetupData
		err := rows.Scan(&s.ID, &s.Repo, &s.FromDate, &s.ToDate, &s.CreatedAt, &s.UpdatedAt)
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
	repo = strings.ToLower(repo)
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
	repo = strings.ToLower(repo)
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
