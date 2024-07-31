package v1

import (
	"github.com/gorilla/mux"
)

func V1routes(r *mux.Router) {
	r.PathPrefix("/api/v1").Subrouter()

}
