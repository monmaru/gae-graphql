package gql

import "github.com/graphql-go/graphql"

func newRootQuery(r resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: newUserType(r),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.queryUser,
			},
			"blogs": makeListField(makeNodeListType("BlogList", newBlogType()), r.queryBlogs),
		},
	})
}
