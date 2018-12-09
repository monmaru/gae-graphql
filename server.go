package main

import (
	"log"
	"net/http"
	"os"

	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/application/usecase"
	"github.com/monmaru/gae-graphql/infrastructure/datastore"
	"github.com/monmaru/gae-graphql/interfaces/router"
)

func main() {
	projID := os.Getenv("PROJECT_ID")
	ud := &datastore.UserDatastore{ProjID: projID}
	bd := &datastore.BlogDatastore{ProjID: projID}
	schema, _ := gql.NewSchema(ud, bd)
	usecase := &usecase.GraphQLInteractor{Schema: schema}
	router := router.New(&schema, usecase)
	http.Handle("/", router)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
