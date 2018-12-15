package gql

import (
	"errors"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
	"github.com/monmaru/gae-graphql/library/profile"
)

type resolver interface {
	queryUser(params graphql.ResolveParams) (interface{}, error)
	queryBlogs(params graphql.ResolveParams) (interface{}, error)
	queryBlogsByUser(params graphql.ResolveParams) (interface{}, error)
	createUser(params graphql.ResolveParams) (interface{}, error)
	createBlog(params graphql.ResolveParams) (interface{}, error)
	createUsersBatch(params graphql.ResolveParams) (interface{}, error)
	createBlogsBatch(params graphql.ResolveParams) (interface{}, error)
}

type userKey struct {
	key  string
	user *model.User
}

func newUserKey(key string, user *model.User) *userKey {
	return &userKey{
		key:  key,
		user: user,
	}
}

func (uk *userKey) String() string {
	return uk.key
}

func (rk *userKey) Raw() interface{} {
	return rk.user
}

type blogKey struct {
	key  string
	blog *model.Blog
}

func newBlogKey(key string, blog *model.Blog) *blogKey {
	return &blogKey{
		key:  key,
		blog: blog,
	}
}

func (bk *blogKey) String() string {
	return bk.key
}

func (bk *blogKey) Raw() interface{} {
	return bk.blog
}

type graphQLResolver struct {
	ur repository.UserRepository
	br repository.BlogRepository
}

func newResolver(ur repository.UserRepository, br repository.BlogRepository) resolver {
	return &graphQLResolver{
		ur: ur,
		br: br,
	}
}

func (r *graphQLResolver) createUser(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.createUser]")
	ctx := params.Context
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	user := &model.User{
		Name:  name,
		EMail: email,
	}
	return r.ur.Create(ctx, user)
}

func (r *graphQLResolver) queryUser(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.queryUser]")
	ctx := params.Context
	if strID, ok := params.Args["id"].(string); ok {
		return r.ur.Get(ctx, strID)
	}
	return model.User{}, nil
}

func (r *graphQLResolver) createBlog(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.createBlog]")
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
	return r.br.Create(ctx, blog)
}

func (r *graphQLResolver) queryBlogs(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.queryBlogs]")
	ctx := params.Context
	query := r.br.NewQuery()
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
	defer profile.Duration(time.Now(), "[graphQLResolver.queryBlogsByUser]")
	ctx := params.Context
	query := r.br.NewQuery()
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

func (r *graphQLResolver) createUsersBatch(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.createUsers]")
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	user := &model.User{
		Name:  name,
		EMail: email,
	}

	key := newUserKey(user.Name+user.EMail, user)
	v := params.Context.Value(CreateUsersKey)
	loader, ok := v.(*dataloader.Loader)
	if !ok {
		return nil, errors.New("loader is empty")
	}

	thunk := loader.Load(params.Context, key)
	return func() (interface{}, error) {
		return thunk()
	}, nil
}

func (r *graphQLResolver) createBlogsBatch(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.createBlogs]")
	title, _ := params.Args["title"].(string)
	content, _ := params.Args["content"].(string)
	userID, _ := params.Args["userId"].(string)
	blog := &model.Blog{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	key := newBlogKey(blog.UserID+blog.Title, blog)
	v := params.Context.Value(CreateBlogsKey)
	loader, ok := v.(*dataloader.Loader)
	if !ok {
		return nil, errors.New("loader is empty")
	}

	thunk := loader.Load(params.Context, key)
	return func() (interface{}, error) {
		return thunk()
	}, nil
}
