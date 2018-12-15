package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/domain/repository"
)

func NewSchema(ur repository.UserRepository, br repository.BlogRepository) (graphql.Schema, error) {
	resolver := newResolver(ur, br)
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    newQuery(resolver),
		Mutation: newMutation(resolver),
	})
}
