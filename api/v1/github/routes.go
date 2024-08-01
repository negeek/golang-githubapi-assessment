package github

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	router := r.PathPrefix("/github").Subrouter()
	router.HandleFunc("/setup/", Setup).Methods("POST")
	router.HandleFunc("/top-n-commit-authors/", TopNCommitAuthors).Methods("POST")
	router.HandleFunc("/repo-commits/", RepoCommits).Methods("POST")

}
