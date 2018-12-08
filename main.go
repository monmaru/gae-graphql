package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
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
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	r := mux.NewRouter()
	r.Path("/api/graphql").Handler(controller.New(schema)).Methods("POST")
	r.Path("/graphiql").Handler(h)
	http.Handle("/", r)
	appengine.Main()
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
