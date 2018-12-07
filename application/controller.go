package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/infrastructure"
	"google.golang.org/appengine"
)

func GraphQLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSON(
			w,
			errorResponse{Message: "Invalid request body"},
			http.StatusBadRequest)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema: gql.NewSchema(
			&infrastructure.UserDatastore{},
			&infrastructure.BlogDatastore{}),
		RequestString: string(body),
		Context:       ctx,
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
