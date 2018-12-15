package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/application/usecase"
	"github.com/monmaru/gae-graphql/domain/repository"
	"github.com/monmaru/gae-graphql/interfaces/handler"
	"github.com/monmaru/gae-graphql/interfaces/middleware"
)

func Build(ur repository.UserRepository, br repository.BlogRepository) (http.Handler, error) {
	router := mux.NewRouter()
	schema, err := gql.NewSchema(ur, br)
	if err != nil {
		return nil, err
	}

	injector := middleware.NewInjector(ur, br)
	usecase := usecase.NewGraphQLUsecae(schema)

	router.Path("/ping").HandlerFunc(handler.Pong).Methods(http.MethodGet)
	router.Path("/api/graphql").
		Handler(middleware.Timetrack(injector.Inject(handler.API(usecase)))).
		Methods(http.MethodPost)
	router.Path("/graphiql").
		Handler(middleware.Timetrack(injector.Inject(handler.GraphiQL(&schema))))
	router.Path("/playground").
		Handler(middleware.Timetrack(injector.Inject(handler.Playground(&schema))))
	return router, nil
}
