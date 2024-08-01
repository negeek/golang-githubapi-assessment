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

	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func ParseCommitData(data []map[string]interface{}, repo string) (commits []githubModels.Commit, err error) {
	for _, datum := range data {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic occurred: %v", r)
				}
			}()

			commitData := datum["commit"].(map[string]interface{})
			sha := commitData["sha"].(string)
			url := commitData["url"].(string)

			authorData := commitData["author"].(map[string]interface{})
			authorName := authorData["name"].(string)
			authorEmail := authorData["email"].(string)
			message := commitData["message"].(string)
			dateStr := authorData["date"].(string)

			date, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				log.Printf("invalid date format: %v", err)
				return
			}

			commit := githubModels.Commit{
				Repo:        repo,
				SHA:         sha,
				URL:         url,
				AuthorName:  authorName,
				AuthorEmail: authorEmail,
				Message:     message,
				Date:        date,
			}
			commits = append(commits, commit)
		}()
	}

	return commits, nil
}

func ParseRepoData(data map[string]interface{}, owner string) (repo *githubModels.Repository, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	createdAt, err := time.Parse(time.RFC3339, data["created_at"].(string))
	if err != nil {
		return nil, errors.New("invalid created_at format")
	}

	updatedAt, err := time.Parse(time.RFC3339, data["updated_at"].(string))
	if err != nil {
		return nil, errors.New("invalid updated_at format")
	}

	repo = &githubModels.Repository{
		Name:            data["name"].(string),
		Description:     data["description"].(string),
		URL:             data["url"].(string),
		Language:        data["language"].(string),
		ForksCount:      data["forks_count"].(int),
		StarsCount:      data["stars_count"].(int),
		OpenIssuesCount: data["open_issues_count"].(int),
		WatchersCount:   data["watchers_count"].(int),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	return repo, nil
}

func FetchSaveCommits(config githubModels.SetupData) {
	var (
		url           string
		urlWithParams string
		err           error
		req           *http.Request
		resp          *http.Response
		data          []map[string]interface{}
		respBody      []byte
		commits       []githubModels.Commit
		queryParams   = make(map[string]string)
	)

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
		queryParams["page"] = strconv.Itoa(page)
		urlWithParams, err = utils.AddQueryParams(url, queryParams)
		if err != nil {
			log.Println(err)
			break
		}

		req, err = http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			log.Println(err)
			break
		}
		req.Header.Set("Accept", "application/vnd.github+json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err)
			break
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("request failed with status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			break
		}

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			break
		}

		resp.Body.Close()

		err = json.Unmarshal(respBody, &data)
		if err != nil {
			log.Println(err)
			break
		}

		if len(data) == 0 {
			break
		}

		commits, err = ParseCommitData(data, config.Repo)
		if err != nil {
			log.Println(err)
			continue
		}

		githubModels.CreateCommits(commits)

		page++
	}
	log.Printf("All commits saved for repo: %s\t and owner:%s", config.Repo, config.Owner)
}

func FetchSaveRepo(config githubModels.SetupData) error {
	var (
		url      string
		err      error
		req      *http.Request
		resp     *http.Response
		data     map[string]interface{}
		respBody []byte
		repo     = &githubModels.Repository{}
	)

	url = fmt.Sprintf(RepoUrl, config.Owner, config.Repo)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to initiate github request")
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("unable to make github request")
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed with status code: %d\n", resp.StatusCode)
		resp.Body.Close()
		log.Println(err)
		return errors.New("github request failed")
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("unable to process github repo response")
	}
	resp.Body.Close()

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Println(err)
		return errors.New("unable to parse json data")
	}

	repo, err = ParseRepoData(data, config.Owner)
	if err != nil {
		log.Println(err)
		return err
	}

	err = repo.Create()
	if err != nil {
		log.Println(err)
		return errors.New("error saving repo data")
	}

	return nil

}
