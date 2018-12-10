package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/monmaru/gae-graphql/application/gql"
)

func main() {
	schema, err := gql.NewSchema(nil, nil)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: testutil.IntrospectionQuery,
	})
	if result.HasErrors() {
		log.Fatalf("ERROR introspecting schema: %v", result.Errors)
		return
	}

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	err = ioutil.WriteFile("./schema.json", b, os.ModePerm)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
