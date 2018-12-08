package controller

import (
	"context"
	"encoding/json"
	"fmt"
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
		writeJSON(
			w,
			errorResponse{Message: "Invalid request body"},
			http.StatusBadRequest)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        c.schema,
		RequestString: string(body),
		Context:       context.Background(),
	})

	if len(result.Errors) > 0 {
		writeJSON(
			w,
			errorResponse{Message: fmt.Sprintf("%+v", result.Errors)},
			http.StatusBadRequest)
		return
	}
	writeJSON(w, result, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

type errorResponse struct {
	Message string `json:"message"`
}
