package infrastructure

import (
	"context"
	"errors"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
)

type UserDatastore struct {
	client *datastore.Client
}

func NewUserDatastore(projID string) (*UserDatastore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projID)
	if err != nil {
		return nil, err
	}
	return &UserDatastore{client: client}, nil
}

func (u *UserDatastore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	key := datastore.IncompleteKey("User", nil)
	generatedKey, err := u.client.Put(ctx, key, user)
	if err != nil {
		return nil, err
	}

	user.ID = strconv.FormatInt(generatedKey.ID, 10)
	return user, nil
}

func (u *UserDatastore) GetUser(ctx context.Context, strID string) (*model.User, error) {
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid id")
	}

	user := &model.User{ID: strID}
	key := datastore.IDKey("User", id, nil)
	if err := u.client.Get(ctx, key, user); err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

type BlogDatastore struct {
	client *datastore.Client
}

func NewBlogDatastore(projID string) (*BlogDatastore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projID)
	if err != nil {
		return nil, err
	}
	return &BlogDatastore{client: client}, nil
}

func (b *BlogDatastore) CreateBlog(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	key := datastore.IncompleteKey("Blog", nil)
	generatedKey, err := b.client.Put(ctx, key, blog)
	if err != nil {
		return nil, err
	}

	blog.ID = strconv.FormatInt(generatedKey.ID, 10)
	return blog, nil
}

func (b *BlogDatastore) NewQuery() repository.Query {
	return &QueryImpl{
		query:  datastore.NewQuery("Blog"),
		client: b.client,
	}
}

type QueryImpl struct {
	query  *datastore.Query
	client *datastore.Client
}

func (q *QueryImpl) Limit(limit int) repository.Query {
	q.query = q.query.Limit(limit)
	return q
}

func (q *QueryImpl) Offset(offset int) repository.Query {
	q.query = q.query.Offset(offset)
	return q
}
func (q *QueryImpl) Order(filedName string) repository.Query {
	q.query = q.query.Order(filedName)
	return q
}

func (q *QueryImpl) Filter(filterStr string, value interface{}) repository.Query {
	q.query = q.query.Filter(filterStr, value)
	return q
}

func (q *QueryImpl) GetAll(ctx context.Context) (*model.BlogList, error) {
	var result model.BlogList
	keys, err := q.client.GetAll(ctx, q.query, &result.Nodes)
	if err != nil {
		return &result, err
	}

	for i, key := range keys {
		result.Nodes[i].ID = strconv.FormatInt(key.ID, 10)
	}
	result.TotalCount = len(result.Nodes)
	return &result, nil
}
