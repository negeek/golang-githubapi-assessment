package github

import (
	"testing"
	"time"

	"github.com/negeek/golang-githubapi-assessment/data/v1/github"
)

var (
	MockCommitData = []map[string]interface{}{
		{
			"commit": map[string]interface{}{
				"url": "https://api.github.com/repos/octocat/hello-world/commits/762941318ee16e59dabbacb1b4049eec22f0d303",
				"author": map[string]interface{}{
					"name":  "monalisa octocat",
					"email": "support@github.com",
					"date":  "2011-04-14T16:00:49Z",
				},
				"message": "Fix all the bugs",
			},
			"sha": "762941318ee16e59dabbacb1b4049eec22f0d303",
		},
	}

	MockRepoData = map[string]interface{}{
		"full_name":         "octocat/hello-world",
		"description":       "This your first repo!",
		"url":               "https://api.github.com/repos/octocat/hello-world",
		"language":          "Go",
		"forks_count":       9.0,
		"stargazers_count":  80.0,
		"open_issues_count": 0.0,
		"watchers_count":    80.0,
		"created_at":        "2011-04-14T16:00:49Z",
		"updated_at":        "2011-04-14T16:00:49Z",
	}
)

func TestParseCommitData(t *testing.T) {
	commits, err := ParseCommitData(MockCommitData, "octocat/hello-world")
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	if len(commits) != 1 {
		t.Errorf("Expected 1 commit, got %d", len(commits))
	}

	expectedCommit := github.Commit{
		Repo:        "octocat/hello-world",
		SHA:         "762941318ee16e59dabbacb1b4049eec22f0d303",
		URL:         "https://api.github.com/repos/octocat/hello-world/commits/762941318ee16e59dabbacb1b4049eec22f0d303",
		AuthorName:  "monalisa octocat",
		AuthorEmail: "support@github.com",
		Message:     "Fix all the bugs",
		Date:        time.Date(2011, time.April, 14, 16, 0, 49, 0, time.UTC),
	}

	commit := commits[0]

	if commit.Repo != expectedCommit.Repo {
		t.Errorf("Expected Repo %s, got %s", expectedCommit.Repo, commit.Repo)
	}
	if commit.SHA != expectedCommit.SHA {
		t.Errorf("Expected SHA %s, got %s", expectedCommit.SHA, commit.SHA)
	}
	if commit.URL != expectedCommit.URL {
		t.Errorf("Expected URL %s, got %s", expectedCommit.URL, commit.URL)
	}
	if commit.AuthorName != expectedCommit.AuthorName {
		t.Errorf("Expected AuthorName %s, got %s", expectedCommit.AuthorName, commit.AuthorName)
	}
	if commit.AuthorEmail != expectedCommit.AuthorEmail {
		t.Errorf("Expected AuthorEmail %s, got %s", expectedCommit.AuthorEmail, commit.AuthorEmail)
	}
	if commit.Message != expectedCommit.Message {
		t.Errorf("Expected Message %s, got %s", expectedCommit.Message, commit.Message)
	}
	if !commit.Date.Equal(expectedCommit.Date) {
		t.Errorf("Expected Date %v, got %v", expectedCommit.Date, commit.Date)
	}
}

func TestParseRepoData(t *testing.T) {
	repo, err := ParseRepoData(MockRepoData)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}

	expectedRepo := &github.Repository{
		Name:            "octocat/hello-world",
		Description:     "This your first repo!",
		URL:             "https://api.github.com/repos/octocat/hello-world",
		Language:        "Go",
		ForksCount:      9.0,
		StarsCount:      80.0,
		OpenIssuesCount: 0.0,
		WatchersCount:   80.0,
		CreatedAt:       time.Date(2011, time.April, 14, 16, 0, 49, 0, time.UTC),
		UpdatedAt:       time.Date(2011, time.April, 14, 16, 0, 49, 0, time.UTC),
	}

	if repo.Name != expectedRepo.Name {
		t.Errorf("Expected Name %s, got %s", expectedRepo.Name, repo.Name)
	}
	if repo.Description != expectedRepo.Description {
		t.Errorf("Expected Description %s, got %s", expectedRepo.Description, repo.Description)
	}
	if repo.URL != expectedRepo.URL {
		t.Errorf("Expected URL %s, got %s", expectedRepo.URL, repo.URL)
	}
	if repo.Language != expectedRepo.Language {
		t.Errorf("Expected Language %s, got %s", expectedRepo.Language, repo.Language)
	}
	if repo.ForksCount != expectedRepo.ForksCount {
		t.Errorf("Expected ForksCount %f, got %f", expectedRepo.ForksCount, repo.ForksCount)
	}
	if repo.StarsCount != expectedRepo.StarsCount {
		t.Errorf("Expected StarsCount %f, got %f", expectedRepo.StarsCount, repo.StarsCount)
	}
	if repo.OpenIssuesCount != expectedRepo.OpenIssuesCount {
		t.Errorf("Expected OpenIssuesCount %f, got %f", expectedRepo.OpenIssuesCount, repo.OpenIssuesCount)
	}
	if repo.WatchersCount != expectedRepo.WatchersCount {
		t.Errorf("Expected WatchersCount %f, got %f", expectedRepo.WatchersCount, repo.WatchersCount)
	}
	if !repo.CreatedAt.Equal(expectedRepo.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", expectedRepo.CreatedAt, repo.CreatedAt)
	}
	if !repo.UpdatedAt.Equal(expectedRepo.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", expectedRepo.UpdatedAt, repo.UpdatedAt)
	}
}
