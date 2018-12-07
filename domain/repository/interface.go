package repository

import (
	"context"

	"github.com/monmaru/gae-graphql/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, strID string) (*model.User, error)
}

type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *model.Blog) (*model.Blog, error)
	NewQuery() Query
}

type Query interface {
	Limit(limit int) Query
	Offset(offset int) Query
	Order(filedName string) Query
	Filter(filterStr string, value interface{}) Query
	GetAll(ctx context.Context) (*model.BlogList, error)
}
