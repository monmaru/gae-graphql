package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/domain/repository"
)

// NewSchema ...
func NewSchema(userRepo repository.UserRepository, blogRepo repository.BlogRepository) graphql.Schema {
	resolver := &graphQLResolver{
		userRepo: userRepo,
		blogRepo: blogRepo,
	}
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    newRootQuery(resolver),
		Mutation: newRootMutation(resolver),
	})
	return schema
}
