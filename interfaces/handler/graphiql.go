package handler

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func GraphiQL(schema *graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: false,
	})
}
