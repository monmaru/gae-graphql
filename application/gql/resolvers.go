package gql

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
)

type resolver interface {
	queryUser(params graphql.ResolveParams) (interface{}, error)
	queryBlogs(params graphql.ResolveParams) (interface{}, error)
	queryBlogsByUser(params graphql.ResolveParams) (interface{}, error)
	createUser(params graphql.ResolveParams) (interface{}, error)
	createBlog(params graphql.ResolveParams) (interface{}, error)
}

type graphQLResolver struct {
	userRepo repository.UserRepository
	blogRepo repository.BlogRepository
}

func newResolver(userRepo repository.UserRepository, blogRepo repository.BlogRepository) resolver {
	return &graphQLResolver{
		userRepo: userRepo,
		blogRepo: blogRepo,
	}
}

func (r *graphQLResolver) createUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	user := &model.User{
		Name:  name,
		EMail: email,
	}
	return r.userRepo.CreateUser(ctx, user)
}

func (r *graphQLResolver) queryUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	if strID, ok := params.Args["id"].(string); ok {
		return r.userRepo.GetUser(ctx, strID)
	}
	return model.User{}, nil
}

func (r *graphQLResolver) createBlog(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	title, _ := params.Args["title"].(string)
	content, _ := params.Args["content"].(string)
	userID, _ := params.Args["userId"].(string)
	blog := &model.Blog{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	return r.blogRepo.CreateBlog(ctx, blog)
}

func (r *graphQLResolver) queryBlogs(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := r.blogRepo.NewQuery()
	query = query.Order("-CreatedAt")
	if limit, ok := params.Args["limit"].(int); ok {
		query = query.Limit(limit)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		query = query.Offset(offset)
	}
	return query.GetAll(ctx)
}

func (r *graphQLResolver) queryBlogsByUser(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	query := r.blogRepo.NewQuery()
	query = query.Order("-CreatedAt")
	if limit, ok := params.Args["limit"].(int); ok {
		query = query.Limit(limit)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		query = query.Offset(offset)
	}
	if user, ok := params.Source.(*model.User); ok {
		query = query.Filter("UserID =", user.ID)
	}
	return query.GetAll(ctx)
}