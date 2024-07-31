package github

import (
	"context"

	"github.com/negeek/golang-githubapi-assessment/db"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

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
	query := "INSERT INTO commits (id, owner, name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	_, err := db.PostgreSQLDB.Exec(context.Background(), query, r.ID, r.Owner, r.Name, r.Description, r.URL, r.Language, r.ForksCount, r.StarsCount, r.OpenIssuesCount, r.WatchersCount, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
