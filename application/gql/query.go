package gql

import "github.com/graphql-go/graphql"

func newQuery(r resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        newUserType(r),
				Description: "Look up a user by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.queryUser,
			},
			"blogs": makeListField(makeNodeListType("BlogList", newBlogType()), r.queryBlogs),
		},
	})
}
