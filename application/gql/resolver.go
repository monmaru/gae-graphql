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
	getUser(params graphql.ResolveParams) (interface{}, error)
	queryBlogs(params graphql.ResolveParams) (interface{}, error)
	queryBlogsByUser(params graphql.ResolveParams) (interface{}, error)
	createUser(params graphql.ResolveParams) (interface{}, error)
	createBlog(params graphql.ResolveParams) (interface{}, error)
	getUsersBatch(params graphql.ResolveParams) (interface{}, error)
	createUsersBatch(params graphql.ResolveParams) (interface{}, error)
	createBlogsBatch(params graphql.ResolveParams) (interface{}, error)
}

type getUserKey struct {
	strID string
}

func newGetUserKey(strID string) *getUserKey {
	return &getUserKey{strID: strID}
}

func (uk *getUserKey) String() string {
	return uk.strID
}

func (uk *getUserKey) Raw() interface{} {
	return uk.strID
}

type createUserKey struct {
	key  string
	user *model.User
}

func newCreateUserKey(key string, user *model.User) *createUserKey {
	return &createUserKey{
		key:  key,
		user: user,
	}
}

func (uk *createUserKey) String() string {
	return uk.key
}

func (uk *createUserKey) Raw() interface{} {
	return uk.user
}

type createBlogKey struct {
	key  string
	blog *model.Blog
}

func newCreateBlogKey(key string, blog *model.Blog) *createBlogKey {
	return &createBlogKey{
		key:  key,
		blog: blog,
	}
}

func (bk *createBlogKey) String() string {
	return bk.key
}

func (bk *createBlogKey) Raw() interface{} {
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

func (r *graphQLResolver) getUser(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.getUser]")
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

func (r *graphQLResolver) getUsersBatch(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.getUsersBatch]")
	strID, ok := params.Args["id"].(string)
	if !ok {
		return nil, errors.New("invalid id")
	}

	key := newGetUserKey(strID)
	v := params.Context.Value(GetUsersKey)
	loader, ok := v.(*dataloader.Loader)
	if !ok {
		return nil, errors.New("loader is empty")
	}

	thunk := loader.Load(params.Context, key)
	return func() (interface{}, error) {
		return thunk()
	}, nil
}

func (r *graphQLResolver) createUsersBatch(params graphql.ResolveParams) (interface{}, error) {
	defer profile.Duration(time.Now(), "[graphQLResolver.createUsersBatch]")
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	user := &model.User{
		Name:  name,
		EMail: email,
	}

	key := newCreateUserKey(user.Name+user.EMail, user)
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
	defer profile.Duration(time.Now(), "[graphQLResolver.createBlogsBatch]")
	title, _ := params.Args["title"].(string)
	content, _ := params.Args["content"].(string)
	userID, _ := params.Args["userId"].(string)
	blog := &model.Blog{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	key := newCreateBlogKey(blog.UserID+blog.Title, blog)
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
