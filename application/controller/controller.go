package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
)

type GraphQLController struct {
	schema graphql.Schema
}

func New(schema graphql.Schema) *GraphQLController {
	return &GraphQLController{schema: schema}
}

func (c *GraphQLController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        c.schema,
		RequestString: string(body),
		Context:       context.Background(),
	})

	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
