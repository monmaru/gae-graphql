package middleware

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/domain/repository"
)

type Injector struct {
	ur repository.UserRepository
	br repository.BlogRepository
}

func NewInjector(ur repository.UserRepository, br repository.BlogRepository) *Injector {
	return &Injector{
		ur: ur,
		br: br,
	}
}

func (i *Injector) Inject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(
			ctx,
			gql.CreateUsersKey,
			dataloader.NewBatchedLoader(gql.CreateUsersBatchFunc(i.ur)))
		ctx = context.WithValue(
			ctx,
			gql.CreateBlogsKey,
			dataloader.NewBatchedLoader(gql.CreateBlogsBatchFunc(i.br)))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
