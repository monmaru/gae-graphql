package datastore

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
)

type BlogDatastore struct {
	client *datastore.Client
	kind   string
}

func NewBlogDatastore(client *datastore.Client) *BlogDatastore {
	return &BlogDatastore{client: client, kind: "Blog"}
}

func (b *BlogDatastore) Create(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	key := datastore.IncompleteKey(b.kind, nil)
	generatedKey, err := b.client.Put(ctx, key, blog)
	if err != nil {
		return nil, err
	}

	blog.ID = strconv.FormatInt(generatedKey.ID, 10)
	return blog, nil
}

func (b *BlogDatastore) CreateMulti(ctx context.Context, blogs []*model.Blog) ([]*model.Blog, error) {
	var keys []*datastore.Key
	for i := 0; i < len(blogs); i++ {
		key := datastore.IncompleteKey(b.kind, nil)
		keys = append(keys, key)
	}

	generatedKeys, err := b.client.PutMulti(ctx, keys, blogs)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(blogs); i++ {
		blogs[i].ID = strconv.FormatInt(generatedKeys[i].ID, 10)
	}
	return blogs, nil
}

func (b *BlogDatastore) NewQuery() repository.Query {
	return &QueryImpl{
		client: b.client,
		query:  datastore.NewQuery(b.kind),
	}
}

type QueryImpl struct {
	client *datastore.Client
	query  *datastore.Query
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
