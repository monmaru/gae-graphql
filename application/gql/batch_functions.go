package gql

import (
	"context"
	"errors"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
	"github.com/monmaru/gae-graphql/library/log"
)

type BatchKey string

const (
	GetUsersKey    BatchKey = "GetUsersBatchKey"
	CreateUsersKey BatchKey = "CreateUsersBatchKey"
	CreateBlogsKey BatchKey = "CreateBlogsBatchKey"
)

func GetUsersBatchFunc(ur repository.UserRepository) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		defer log.Duration(ctx, time.Now(), "[GetUsersBatchFunc]")
		var strIDs []string
		for _, key := range keys {
			strID, ok := key.Raw().(string)
			if !ok {
				return handleError(ctx, errors.New("Invalid key value"))
			}
			strIDs = append(strIDs, strID)
		}

		users, err := ur.GetMulti(ctx, strIDs)
		if err != nil {
			log.Errorf(ctx, err.Error())
			return handleError(ctx, err)
		}

		var results []*dataloader.Result
		for _, user := range users {
			result := dataloader.Result{
				Data:  user,
				Error: nil,
			}
			results = append(results, &result)
		}

		log.Infof(ctx, "[GetUsersBatchFunc] batch size: %d", len(results))
		return results
	}
}

func CreateUsersBatchFunc(ur repository.UserRepository) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		defer log.Duration(ctx, time.Now(), "[CreateUsersBatchFn]")
		var users []*model.User
		for _, key := range keys {
			u, ok := key.Raw().(*model.User)
			if !ok {
				return handleError(ctx, errors.New("Invalid key value"))
			}
			users = append(users, u)
		}

		savedUsers, err := ur.CreateMulti(ctx, users)
		if err != nil {
			return handleError(ctx, err)
		}

		var results []*dataloader.Result
		for _, user := range savedUsers {
			result := dataloader.Result{
				Data:  user,
				Error: nil,
			}
			results = append(results, &result)
		}

		log.Infof(ctx, "[CreateUsersBatchFn] batch size: %d", len(results))
		return results
	}
}

func CreateBlogsBatchFunc(br repository.BlogRepository) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		defer log.Duration(ctx, time.Now(), "[CreateBlogsBatchFn]")
		var blogs []*model.Blog
		for _, key := range keys {
			b, ok := key.Raw().(*model.Blog)
			if !ok {
				return handleError(ctx, errors.New("Invalid key value"))
			}
			blogs = append(blogs, b)
		}

		savedBlogs, err := br.CreateMulti(ctx, blogs)
		if err != nil {
			return handleError(ctx, err)
		}

		var results []*dataloader.Result
		for _, blog := range savedBlogs {
			result := dataloader.Result{
				Data:  blog,
				Error: nil,
			}
			results = append(results, &result)
		}

		log.Infof(ctx, "[CreateBlogsBatchFn] batch size: %d", len(results))
		return results
	}
}

func handleError(ctx context.Context, err error) []*dataloader.Result {
	log.Errorf(ctx, err.Error())
	var results []*dataloader.Result
	var result dataloader.Result
	result.Error = err
	results = append(results, &result)
	return results
}
