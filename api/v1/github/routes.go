package github

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	r.PathPrefix("/github").Subrouter()
}
