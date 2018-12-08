package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monmaru/gae-graphql/application"
	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()
	r.Methods("POST").Path("/graphql").HandlerFunc(application.GraphQLHandler)
	http.Handle("/", r)
	appengine.Main()
}
