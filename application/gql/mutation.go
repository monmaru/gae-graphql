package gql

import "github.com/graphql-go/graphql"

func newRootMutation(r resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: newCreateUserInputType(r),
				Args: graphql.FieldConfigArgument{
					"name":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.createUser,
			},
			"createBlog": &graphql.Field{
				Type: newCreateBlogInputType(),
				Args: graphql.FieldConfigArgument{
					"userId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"title":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.createBlog,
			},
		},
	})
}
