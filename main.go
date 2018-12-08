package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/monmaru/gae-graphql/application/controller"
	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/infrastructure"
	"google.golang.org/appengine"
)

func main() {
	projID := os.Getenv("PROJECT_ID")
	ud, err := infrastructure.NewUserDatastore(projID)
	exitIfError(err)

	bd, err := infrastructure.NewBlogDatastore(projID)
	exitIfError(err)

	schema, _ := gql.NewSchema(ud, bd)

	r := mux.NewRouter()
	r.Methods("POST").Path("/graphql").Handler(controller.New(schema))
	http.Handle("/", r)
	appengine.Main()
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
