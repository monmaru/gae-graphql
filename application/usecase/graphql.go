package usecase

import (
	"context"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/library/profile"
)

type GraphQLUsecase interface {
	Do(ctx context.Context, body string) interface{}
}

type GraphQLInteractor struct {
	Schema graphql.Schema
}

func (i *GraphQLInteractor) Do(ctx context.Context, body string) interface{} {
	defer profile.Duration(time.Now(), "GraphQLInteractor.Do")
	return graphql.Do(graphql.Params{
		Schema:        i.Schema,
		RequestString: string(body),
		Context:       ctx,
	})
}
