package github

import (
	"log"
	"time"

	githubFuncs "github.com/negeek/golang-githubapi-assessment/api/v1/github"
	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
)

func CommitCron() {
	var (
		setups = []githubModels.SetupData{}
		err    error
		commit = &githubModels.Commit{}
		setup  = githubModels.SetupData{}
	)

	setups, err = githubModels.Get_all_setup_data()
	if err != nil {
		log.Println(err)
		return
	}

	for _, s := range setups {
		commit.Repo = s.Repo
		err = commit.FindLatestRepoCommitByDate()
		if err != nil {
			log.Println((err))
			return
		}
		setup.Repo = s.Repo
		setup.Owner = s.Owner
		setup.FromDate = commit.Date.Add(time.Second)
		setup.ToDate = time.Time{}
		githubFuncs.FetchSaveCommits(setup)
	}
}
