package github

import (
	"net/http"

	"github.com/negeek/golang-githubapi-assessment/utils"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	var (
		payload = SetupData{}
		err     error
	)

	// read  request body
	err = utils.Unmarshall(r.Body, &payload)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, "error reading payload", nil)
		return
	}

	// validate request body
	err = payload.Validate()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// start daemon processes
	go func(data SetupData) {
		FetchSaveRepo(data)
	}(payload)

	go func(data SetupData) {
		FetchSaveCommits(data)
	}(payload)

	utils.JsonResponse(w, true, http.StatusOK, "github setup started successfully", nil)
}
