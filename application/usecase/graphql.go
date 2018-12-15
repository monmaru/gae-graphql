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

func NewGraphQLUsecae(scheme graphql.Schema) GraphQLUsecase {
	return &graphQLInteractor{schema: scheme}
}

type graphQLInteractor struct {
	schema graphql.Schema
}

func (i *graphQLInteractor) Do(ctx context.Context, body string) interface{} {
	defer profile.Duration(time.Now(), "[graphQLInteractor.Do]")
	return graphql.Do(graphql.Params{
		Schema:        i.schema,
		RequestString: string(body),
		Context:       ctx,
	})
}
