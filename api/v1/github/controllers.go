package github

import (
	"net/http"

	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	var (
		requestData = githubModels.SetupData{}
		err         error
	)

	// read  request body
	err = utils.Unmarshall(r.Body, &requestData)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error reading request data", nil)
		return
	}

	// validate request body
	err = requestData.Validate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// start daemon processes
	go func(data githubModels.SetupData) {
		FetchSaveRepo(data)
	}(requestData)

	go func(data githubModels.SetupData) {
		FetchSaveCommits(data)
	}(requestData)

	utils.JsonResponse(w, true, http.StatusOK, "github setup started successfully", nil)
}

func TopNCommitAuthors(w http.ResponseWriter, r *http.Request) {
	var (
		requestData  = TopNCommitAuthorsRequestData{}
		err          error
		responseData []map[string]interface{}
	)

	// read  request body
	err = utils.Unmarshall(r.Body, &requestData)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error reading request data", nil)
		return
	}

	// validate request body
	err = requestData.Validate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	responseData, err = githubModels.GetTopNCommitAuthors(requestData.Repo, requestData.TopN)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusInternalServerError, "error getting topN commit authors", nil)
		return
	}

	utils.JsonResponse(w, true, http.StatusOK, "topN commit authors", responseData)

}

func RepoCommits(w http.ResponseWriter, r *http.Request) {
	var (
		requestData  = RepoCommitsRequestData{}
		err          error
		responseData []map[string]interface{}
	)

	// read  request body
	err = utils.Unmarshall(r.Body, &requestData)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error reading request data", nil)
		return
	}

	// validate request body
	err = requestData.Validate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	responseData, err = githubModels.GetCommitsByRepoName(requestData.Repo)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusInternalServerError, "error getting repo commits", nil)
		return
	}

	utils.JsonResponse(w, true, http.StatusOK, "repo commits", responseData)

}