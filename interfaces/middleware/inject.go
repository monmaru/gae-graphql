package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/monmaru/gae-graphql/application/gql"
	"github.com/monmaru/gae-graphql/domain/repository"
	"github.com/monmaru/gae-graphql/library/log"
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
		ctx := i.setupContext(r)
		defer log.Duration(ctx, time.Now(), fmt.Sprintf("[%s]", r.URL.Path))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (i *Injector) setupContext(r *http.Request) context.Context {
	ctx := r.Context()
	ctx = context.WithValue(
		ctx,
		gql.GetUsersKey,
		dataloader.NewBatchedLoader(gql.GetUsersBatchFunc(i.ur)))
	ctx = context.WithValue(
		ctx,
		gql.CreateUsersKey,
		dataloader.NewBatchedLoader(gql.CreateUsersBatchFunc(i.ur)))
	ctx = context.WithValue(
		ctx,
		gql.CreateBlogsKey,
		dataloader.NewBatchedLoader(gql.CreateBlogsBatchFunc(i.br)))
	return log.WithContext(ctx, r)
}
