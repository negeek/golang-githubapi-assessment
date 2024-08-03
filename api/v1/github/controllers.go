package github

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
	"github.com/negeek/golang-githubapi-assessment/utils"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	rawData := make(map[string]string)
	err := utils.Unmarshall(r.Body, &rawData)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error reading request data", nil)
		return
	}

	requestData := githubModels.SetupData{
		Repo:     utils.ExtractStringField(rawData, "repo"),
		FromDate: time.Time{}, // default value
		ToDate:   time.Now(),  // default value
	}

	requestData.FromDate, err = utils.HandleDateField(rawData["from_date"], requestData.FromDate)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "invalid from_date format", nil)
		return
	}

	requestData.ToDate, err = utils.HandleDateField(rawData["to_date"], requestData.ToDate)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "invalid to_date format", nil)
		return
	}

	err = requestData.Validate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// make request to github to get repo metadata
	err = RepoHandler(requestData)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error fetching and saving repo metadata", nil)
		return
	}

	err = requestData.CreateOrUpdate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error saving setup data", nil)
		return
	}

	utils.JsonResponse(w, true, http.StatusOK, "github setup started successfully", nil)
}

func TopNCommitAuthors(w http.ResponseWriter, r *http.Request) {
	var (
		repo         string
		repo_name    string
		owner_name   string
		topN         int
		err          error
		responseData []map[string]interface{}
	)
	pathVars := mux.Vars(r)
	repo_name = pathVars["repo_name"]
	owner_name = pathVars["owner_name"]
	repo = owner_name + "/" + repo_name
	topN, err = strconv.Atoi(pathVars["n"])
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error parsing top n", nil)
		return
	}

	responseData, err = githubModels.GetTopNCommitAuthors(repo, topN)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusInternalServerError, "error getting topN commit authors", nil)
		return
	}

	utils.JsonResponse(w, true, http.StatusOK, "topN commit authors", responseData)

}

func RepoCommits(w http.ResponseWriter, r *http.Request) {
	var (
		repo         string
		repo_name    string
		owner_name   string
		err          error
		responseData []map[string]interface{}
	)

	pathVars := mux.Vars(r)
	repo_name = pathVars["repo_name"]
	owner_name = pathVars["owner_name"]
	repo = owner_name + "/" + repo_name
	responseData, err = githubModels.GetCommitsByRepoName(repo)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusInternalServerError, "error getting repo commits", nil)
		return
	}

	utils.JsonResponse(w, true, http.StatusOK, "repo commits", responseData)

}
