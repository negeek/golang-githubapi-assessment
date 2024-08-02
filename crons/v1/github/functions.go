package github

import (
	"log"
	"time"

	githubFuncs "github.com/negeek/golang-githubapi-assessment/api/v1/github"
	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
)

func CommitCron() {
	/*
		This funtion gets all setup data from db which contains the repo name and owner name.
		It uses this to fetch the commits of each repo. If repo doesn't exist in db,
		it fetches repo detail before fetching commits
	*/
	log.Println("commit cron started")
	var (
		exist  bool
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
			err = githubFuncs.FetchSaveRepo(setup)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		githubFuncs.FetchSaveCommits(setup)
	}
}
