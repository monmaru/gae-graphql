package gql

import "github.com/graphql-go/graphql"

func makeListField(listType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Field {
	return &graphql.Field{
		Type:        listType,
		Description: "Return a list of blogs",
		Resolve:     resolve,
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
	}
}

func makeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes":      &graphql.Field{Type: graphql.NewList(nodeType)},
			"totalCount": &graphql.Field{Type: graphql.Int},
		},
	})
}
