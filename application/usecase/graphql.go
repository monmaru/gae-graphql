package usecase

import (
	"context"

	"github.com/graphql-go/graphql"
)

type GraphQLUsecase interface {
	Do(ctx context.Context, body string) interface{}
}

type GraphQLInteractor struct {
	Schema graphql.Schema
}

func (c *GraphQLInteractor) Do(ctx context.Context, body string) interface{} {
	return graphql.Do(graphql.Params{
		Schema:        c.Schema,
		RequestString: string(body),
		Context:       ctx,
	})
}
