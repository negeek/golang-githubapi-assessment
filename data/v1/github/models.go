package github

import "time"

type Repository struct {
	ID              int       `json:"id"`
	Owner           string    `json:"owner"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	StarsCount      int       `json:"stars_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	WatchersCount   int       `json:"watchers_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Commit struct {
	ID          int        `json:"id"`
	SHA         string     `json:"sha"`
	Repo        string     `json:"repo"` // corresponds to the repo name
	RepoInfo    Repository `json:"repo_info"`
	AuthorName  string     `json:"author_name"`
	AuthorEmail string     `json:"author_email"`
	URL         string     `json:"url"`
	Message     string     `json:"message"`
	Date        time.Time  `json:"date"`
}
