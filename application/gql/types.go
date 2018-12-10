package gql

import "github.com/graphql-go/graphql"

func newUserType(r resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.String},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
			"posts": makeListField(
				makeNodeListType("PostList", newPostType()),
				concurrentResolve(r.queryBlogsByUser)),
		},
	})
}

func newCreateUserInputType(r resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "CreateUserInput",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.String},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	})
}

func newBlogType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Blog",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"userId":    &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: graphql.DateTime},
			"title":     &graphql.Field{Type: graphql.String},
			"content":   &graphql.Field{Type: graphql.String},
		},
	})
}

func newCreateBlogInputType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "CreateBlogInput",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"userId":    &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: graphql.DateTime},
			"title":     &graphql.Field{Type: graphql.String},
			"content":   &graphql.Field{Type: graphql.String},
		},
	})
}

func newPostType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"userId":    &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: graphql.DateTime},
			"title":     &graphql.Field{Type: graphql.String},
			"content":   &graphql.Field{Type: graphql.String},
		},
	})
}
