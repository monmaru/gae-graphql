package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/application/usecase"
	"github.com/monmaru/gae-graphql/interfaces/handler"
	"github.com/monmaru/gae-graphql/interfaces/middleware"
)

func Route(schema *graphql.Schema) http.Handler {
	router := mux.NewRouter()
	usecase := &usecase.GraphQLInteractor{Schema: *schema}
	router.Path("/ping").HandlerFunc(handler.Pong).Methods(http.MethodGet)
	router.Path("/api/graphql").Handler(middleware.Timetrack(handler.API(usecase))).Methods(http.MethodPost)
	router.Path("/graphiql").Handler(middleware.Timetrack(handler.GraphiQL(schema)))
	router.Path("/playground").Handler(middleware.Timetrack(handler.Playground(schema)))
	return router
}
