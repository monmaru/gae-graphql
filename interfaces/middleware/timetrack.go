package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/monmaru/gae-graphql/library/profile"
)

func Timetrack(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer profile.Duration(time.Now(), fmt.Sprintf("[%s]", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
