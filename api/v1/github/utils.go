package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	models "github.com/negeek/golang-githubapi-assessment/data/v1/github"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func ParseCommitData(data []map[string]interface{}, repo string) ([]models.Commit, error) {
	var (
		commits []models.Commit
	)
	for _, datum := range data {
		commitData, ok := datum["commit"].(map[string]interface{})
		if !ok {
			return nil, errors.New("commit data missing or of incorrect type")
		}

		sha, ok := commitData["sha"].(string)
		if !ok {
			return nil, errors.New("commit sha missing or of incorrect type")
		}

		url, ok := commitData["url"].(string)
		if !ok {
			return nil, errors.New("commit url missing or of incorrect type")
		}

		authorData, ok := commitData["author"].(map[string]interface{})
		if !ok {
			return nil, errors.New("author data missing or of incorrect type")
		}

		authorName, ok := authorData["name"].(string)
		if !ok {
			return nil, errors.New("author name missing or of incorrect type")
		}

		authorEmail, ok := authorData["email"].(string)
		if !ok {
			return nil, errors.New("author email missing or of incorrect type")
		}

		message, ok := commitData["message"].(string)
		if !ok {
			return nil, errors.New("commit message missing or of incorrect type")
		}

		dateStr, ok := authorData["date"].(string)
		if !ok {
			return nil, errors.New("author date missing or of incorrect type")
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}

		commit := models.Commit{
			Repo:        repo,
			SHA:         sha,
			URL:         url,
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			Message:     message,
			Date:        date,
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func ParseRepoData(data map[string]interface{}, owner string) (*models.Repository, error) {
	name, ok := data["name"].(string)
	if !ok {
		return nil, errors.New("repo name missing or of incorrect type")
	}

	description, ok := data["description"].(string)
	if !ok {
		return nil, errors.New("repo description missing or of incorrect type")
	}

	url, ok := data["url"].(string)
	if !ok {
		return nil, errors.New("repo url missing or of incorrect type")
	}

	language, ok := data["language"].(string)
	if !ok {
		return nil, errors.New("repo language missing or of incorrect type")
	}

	forks_count, ok := data["forks_count"].(int)
	if !ok {
		return nil, errors.New("repo forks_count missing or of incorrect type")
	}

	stars_count, ok := data["stars_count"].(int)
	if !ok {
		return nil, errors.New("repo stars_count missing or of incorrect type")
	}

	open_issues_count, ok := data["open_issues_count"].(int)
	if !ok {
		return nil, errors.New("repo open_issues_count missing or of incorrect type")
	}

	watchers_count, ok := data["watchers_count"].(int)
	if !ok {
		return nil, errors.New("repo watchers_count missing or of incorrect type")
	}

	created_at_str, ok := data["created_at"].(string)
	if !ok {
		return nil, errors.New("repo created_at missing or of incorrect type")
	}

	created_at, err := time.Parse(time.RFC3339, created_at_str)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at format: %v", err)
	}

	updated_at_str, ok := data["updated_at"].(string)
	if !ok {
		return nil, errors.New("repo updated_at missing or of incorrect type")
	}

	updated_at, err := time.Parse(time.RFC3339, updated_at_str)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at format: %v", err)
	}

	repo := &models.Repository{
		Name:            name,
		Description:     description,
		URL:             url,
		Language:        language,
		ForksCount:      forks_count,
		StarsCount:      stars_count,
		OpenIssuesCount: open_issues_count,
		WatchersCount:   watchers_count,
		CreatedAt:       created_at,
		UpdatedAt:       updated_at,
	}

	return repo, nil

}

func FetchSaveCommits(config SetupData) {
	var (
		url           string
		urlWithParams string
		err           error
		req           *http.Request
		resp          *http.Response
		data          []map[string]interface{}
		respBody      []byte
		commits       []models.Commit
		queryParams   = make(map[string]string)
	)

	queryParams["page"] = "1"
	queryParams["per_page"] = strconv.Itoa(PerPage)
	if !config.FromDate.IsZero() {
		queryParams["since"] = config.FromDate.Format(time.RFC3339)
	}
	if !config.ToDate.IsZero() {
		queryParams["until"] = config.ToDate.Format(time.RFC3339)
	}

	url = fmt.Sprintf(CommitUrl, config.Owner, config.Repo)
	page := 1

	for {
		// Add query parameters to URL
		urlWithParams, err = utils.AddQueryParams(url, queryParams)
		if err != nil {
			log.Println(err)
			return
		}

		req, err = http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Set("Accept", "application/vnd.github+json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request failed with status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			return
		}

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		resp.Body.Close()

		err = json.Unmarshal(respBody, &data)
		if err != nil {
			log.Println(err)
			return
		}

		if len(data) == 0 {
			break
		}

		commits, err = ParseCommitData(data, config.Repo)
		if err != nil {
			log.Println(err)
			return
		}

		err = models.CreateCommits(commits)
		if err != nil {
			log.Println(err)
			return
		}

		page += 1
		queryParams["page"] = strconv.Itoa(page)
	}
	log.Printf("All commits saved for repo: %s\t and owner:%s", config.Repo, config.Owner)
}

func FetchSaveRepo(config SetupData) {
	var (
		url      string
		err      error
		req      *http.Request
		resp     *http.Response
		data     map[string]interface{}
		respBody []byte
		repo     = &models.Repository{}
	)

	url = fmt.Sprintf(RepoUrl, config.Owner, config.Repo)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %d\n", resp.StatusCode)
		resp.Body.Close()
		return
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	resp.Body.Close()

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Println(err)
		return
	}

	repo, err = ParseRepoData(data, config.Owner)
	if err != nil {
		log.Println(err)
		return
	}

	err = repo.Create()
	if err != nil {
		log.Println(err)
		return
	}

}
