package v1

import (
	"github.com/gorilla/mux"
	githubAPI "github.com/negeek/golang-githubapi-assessment/api/v1/github"
)

func V1routes(r *mux.Router) {
	router := r.PathPrefix("/api/v1").Subrouter()
	githubAPI.Routes(router)

}
