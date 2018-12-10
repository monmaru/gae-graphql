package main

import (
	"log"
	"net/http"
	"os"

	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/infrastructure/datastore"
	"github.com/monmaru/gae-graphql/interfaces/router"
)

func main() {
	if err := register(); err != nil {
		log.Fatal(err)
	}
}

func register() error {
	projID := os.Getenv("PROJECT_ID")
	ud := &datastore.UserDatastore{ProjID: projID}
	bd := &datastore.BlogDatastore{ProjID: projID}
	schema, err := gql.NewSchema(ud, bd)
	if err != nil {
		return err
	}
	router := router.Route(&schema)
	http.Handle("/", router)
	return http.ListenAndServe(":8080", router)
}
