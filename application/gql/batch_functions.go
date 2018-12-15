package gql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/graph-gophers/dataloader"

	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
	"github.com/monmaru/gae-graphql/library/profile"
)

type BatchKey string

const (
	CreateUsersKey BatchKey = "CreateUsersBatchKey"
	CreateBlogsKey BatchKey = "CreateBlogsBatchKey"
)

func CreateUsersBatchFunc(ur repository.UserRepository) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		defer profile.Duration(time.Now(), "[CreateUsersBatchFn]")
		var users []*model.User
		for _, key := range keys {
			u, ok := key.Raw().(*model.User)
			if !ok {
				return handleError(errors.New("Invalid key value"))
			}
			users = append(users, u)
		}

		savedUsers, err := ur.CreateMulti(ctx, users)
		if err != nil {
			return handleError(err)
		}

		var results []*dataloader.Result
		for _, user := range savedUsers {
			result := dataloader.Result{
				Data:  user,
				Error: nil,
			}
			results = append(results, &result)
		}

		log.Printf("[CreateUsersBatchFn] batch size: %d", len(results))
		return results
	}
}

func CreateBlogsBatchFunc(br repository.BlogRepository) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		defer profile.Duration(time.Now(), "[CreateBlogsBatchFn]")
		var blogs []*model.Blog
		for _, key := range keys {
			b, ok := key.Raw().(*model.Blog)
			if !ok {
				return handleError(errors.New("Invalid key value"))
			}
			blogs = append(blogs, b)
		}

		savedBlogs, err := br.CreateMulti(ctx, blogs)
		if err != nil {
			log.Println(err.Error())
			return handleError(err)
		}

		var results []*dataloader.Result
		for _, blog := range savedBlogs {
			result := dataloader.Result{
				Data:  blog,
				Error: nil,
			}
			results = append(results, &result)
		}

		log.Printf("[CreateBlogsBatchFn] batch size: %d", len(results))
		return results
	}
}

func handleError(err error) []*dataloader.Result {
	var results []*dataloader.Result
	var result dataloader.Result
	result.Error = err
	results = append(results, &result)
	return results
}
