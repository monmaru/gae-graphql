package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monmaru/gae-graphql/application"
)

func init() {
	r := mux.NewRouter()
	r.Methods("POST").Path("/graphql").HandlerFunc(application.GraphQLHandler)
	http.Handle("/", r)
}
