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
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	projID := os.Getenv("PROJECT_ID")
	ud, err := infrastructure.NewUserDatastore(projID)
	if err != nil {
		return err
	}

	bd, err := infrastructure.NewBlogDatastore(projID)
	if err != nil {
		return err
	}

	schema, _ := gql.NewSchema(ud, bd)
	playground := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	router := mux.NewRouter()
	router.Path("/api/graphql").Handler(controller.New(schema)).Methods(http.MethodPost)
	router.Path("/playground").Handler(playground)
	http.Handle("/", router)
	return http.ListenAndServe(":8080", router)
}
