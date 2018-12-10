package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/application/usecase"
	"github.com/monmaru/gae-graphql/interfaces/handler"
)

func Route(schema *graphql.Schema) http.Handler {
	router := mux.NewRouter()
	usecase := &usecase.GraphQLInteractor{Schema: *schema}
	router.Path("/ping").HandlerFunc(handler.Pong).Methods(http.MethodGet)
	router.Path("/api/graphql").Handler(handler.API(usecase)).Methods(http.MethodPost)
	router.Path("/graphiql").Handler(handler.GraphiQL(schema))
	router.Path("/playground").Handler(handler.Playground(schema))
	return router
}
