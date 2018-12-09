package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/monmaru/gae-graphql/application/usecase"
)

func API(usecase usecase.GraphQLUsecase) http.Handler {
	return &apiHandler{usecase: usecase}
}

type apiHandler struct {
	usecase usecase.GraphQLUsecase
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.usecase.Do(context.Background(), string(body))
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
