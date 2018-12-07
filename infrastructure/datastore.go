package infrastructure

import (
	"context"
	"errors"
	"strconv"

	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
	"google.golang.org/appengine/datastore"
)

type UserDatastore struct{}

func (u *UserDatastore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	key := datastore.NewIncompleteKey(ctx, "User", nil)
	generatedKey, err := datastore.Put(ctx, key, user)
	if err != nil {
		return nil, err
	}

	user.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	return user, nil
}

func (u *UserDatastore) GetUser(ctx context.Context, strID string) (*model.User, error) {
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid id")
	}

	user := &model.User{ID: strID}
	key := datastore.NewKey(ctx, "User", "", id, nil)
	if err := datastore.Get(ctx, key, user); err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

type BlogDatastore struct{}

func (b *BlogDatastore) CreateBlog(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	key := datastore.NewIncompleteKey(ctx, "Blog", nil)
	generatedKey, err := datastore.Put(ctx, key, blog)
	if err != nil {
		return nil, err
	}

	blog.ID = strconv.FormatInt(generatedKey.IntID(), 10)
	return blog, nil
}

func (b *BlogDatastore) NewQuery() repository.Query {
	return &QueryImpl{query: datastore.NewQuery("Blog")}
}

type QueryImpl struct {
	query *datastore.Query
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
	keys, err := q.query.GetAll(ctx, &result.Nodes)
	if err != nil {
		return &result, err
	}

	for i, key := range keys {
		result.Nodes[i].ID = strconv.FormatInt(key.IntID(), 10)
	}
	result.TotalCount = len(result.Nodes)
	return &result, nil
}
