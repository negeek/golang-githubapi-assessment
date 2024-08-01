package github

import "time"

type Repository struct {
	ID              int       `json:"id"`
	Owner           string    `json:"owner"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	Language        string    `json:"language"`
	ForksCount      float64   `json:"forks_count"`
	StarsCount      float64   `json:"stars_count"`
	OpenIssuesCount float64   `json:"open_issues_count"`
	WatchersCount   float64   `json:"watchers_count"`
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

type SetupData struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
