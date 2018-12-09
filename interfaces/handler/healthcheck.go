package handler

import (
	"io"
	"net/http"
)

func Pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "pong")
}
