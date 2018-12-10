package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/domain/repository"
)

// NewSchema ...
func NewSchema(userRepo repository.UserRepository, blogRepo repository.BlogRepository) (graphql.Schema, error) {
	resolver := &graphQLResolver{
		userRepo: userRepo,
		blogRepo: blogRepo,
	}
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    newQuery(resolver),
		Mutation: newMutation(resolver),
	})
}
