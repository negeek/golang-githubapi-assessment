package github

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	router := r.PathPrefix("/github").Subrouter()
	router.HandleFunc("/setup/", Setup).Methods("POST")
	repoRouter := router.PathPrefix("/repo").Subrouter()
	repoRouter.HandleFunc("/{owner_name}/{repo_name}/top/{n:[0-9]+}/commit-authors/", TopNCommitAuthors).Methods("GET")
	repoRouter.HandleFunc("/{owner_name}/{repo_name}/commits/", RepoCommits).Methods("GET")

}
