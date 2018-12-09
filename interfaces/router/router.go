package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/application/usecase"
	"github.com/monmaru/gae-graphql/interfaces/handler"
)

func New(schema *graphql.Schema, usecase usecase.GraphQLUsecase) http.Handler {
	router := mux.NewRouter()
	router.Path("/api/graphql").Handler(handler.API(usecase)).Methods(http.MethodPost)
	router.Path("/playground").Handler(handler.Playground(schema))
	return router
}
