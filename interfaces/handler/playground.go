package handler

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func Playground(schema *graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})
}
