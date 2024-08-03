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

func GithubHandler() {
	/*
		This funtion gets all setup data from db which contains the repo name and owner name.
		It uses this to fetch the commits of each repo. If repo doesn't exist in db,
		it fetches repo detail before fetching commits
	*/
	log.Println("repo commit manager started")
	var (
		exist  bool
		setups = []githubModels.SetupData{}
		err    error
		commit = &githubModels.Commit{}
		setup  = githubModels.SetupData{}
	)

	setups, err = githubModels.GetAllSetUpData()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("gotten all setup data")

	for _, s := range setups {
		commit.Repo = s.Repo
		log.Printf("processing repo: %s", s.Repo)
		exist, err = commit.FindLatestRepoCommitByDate()
		if err != nil {
			log.Println(err)
			continue
		}
		if exist {
			log.Println("commit record for repo exist. Setting FromDate for fetching new commits")
			setup.FromDate = commit.Date
		} else {
			log.Println("no commit record for repo exist.")
			setup.FromDate = time.Time{}
		}
		setup.Repo = s.Repo
		setup.ToDate = time.Time{}

		exist, err = githubModels.FindRepoByName(s.Repo)
		if err != nil {
			log.Println(err)
			continue
		}
		if !exist {
			log.Println("repo does not exist in db. Fetching details")
			err = RepoHandler(setup)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		CommitHandler(setup)
	}
}

func ParseCommitData(data []map[string]interface{}, repo string) (commits []githubModels.Commit, err error) {
	/*
		This function parses the json commit data response of github API into type Commit format.
	*/
	for _, datum := range data {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic occurred: %v, data: %v", r, datum)
				}
			}()

			commitData := datum["commit"].(map[string]interface{})
			sha := datum["sha"].(string)
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

func ParseRepoData(data map[string]interface{}) (repo *githubModels.Repository, err error) {
	/*
		This function parses the json repo data response of github API into type Repository format.
	*/
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v, data: %v", r, data)
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
		Name:            data["full_name"].(string),
		Description:     data["description"].(string),
		URL:             data["url"].(string),
		Language:        data["language"].(string),
		ForksCount:      data["forks_count"].(float64),
		StarsCount:      data["stargazers_count"].(float64),
		OpenIssuesCount: data["open_issues_count"].(float64),
		WatchersCount:   data["watchers_count"].(float64),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	return repo, nil
}

func CommitHandler(config githubModels.SetupData) {
	/*
		This function fetches the repo commits by making request to github and
		then saving the result to db.
	*/
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

	url = fmt.Sprintf(CommitUrl, config.Repo)
	page := 1
	log.Println("starting fetching of commits", queryParams)
	for {
		// Add query parameters to URL
		log.Printf("page %s", strconv.Itoa(page))
		log.Println("parse query params")
		queryParams["page"] = strconv.Itoa(page)
		urlWithParams, err = utils.AddQueryParams(url, queryParams)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("initialize request", urlWithParams)
		req, err = http.NewRequest("GET", urlWithParams, nil)
		if err != nil {
			log.Println(err)
			break
		}
		req.Header.Set("Accept", "application/vnd.github+json")

		log.Println("make request")
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("check response status")
		if resp.StatusCode != http.StatusOK {
			log.Printf("request failed with status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			break
		}

		log.Println("read response body")
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			break
		}

		resp.Body.Close()

		log.Println("parse json response body")
		err = json.Unmarshal(respBody, &data)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("check data length", len(data))
		if len(data) == 0 {
			break
		}

		log.Println("parse commits data")
		commits, err = ParseCommitData(data, config.Repo)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("save commits data")
		githubModels.CreateCommits(commits)

		page++
	}
	log.Printf("All commits saved for repo: %s", config.Repo)
}

func RepoHandler(config githubModels.SetupData) error {
	/*
		This function fetches the repo details by making request to github and
		then saving the result to db.
	*/
	var (
		url      string
		err      error
		req      *http.Request
		resp     *http.Response
		data     map[string]interface{}
		respBody []byte
		repo     = &githubModels.Repository{}
	)

	url = fmt.Sprintf(RepoUrl, config.Repo)
	log.Println("starting fetching repo", url)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to initiate github request")
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	log.Println("make request")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("unable to make github request")
	}

	log.Println("checking status code")
	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed with status code: %d\n", resp.StatusCode)
		resp.Body.Close()
		log.Println(err)
		return errors.New("github request failed")
	}

	log.Println("reading response body")
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return errors.New("unable to process github repo response")
	}
	resp.Body.Close()

	log.Println("parsing json response body")
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Println(err)
		return errors.New("unable to parse json data")
	}

	log.Println("parsing repo data")
	repo, err = ParseRepoData(data)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("saving repo data")
	err = repo.Create()
	if err != nil {
		log.Println(err)
		return errors.New("error saving repo data")
	}

	return nil

}
