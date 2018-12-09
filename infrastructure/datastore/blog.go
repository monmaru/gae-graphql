package datastore

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/monmaru/gae-graphql/domain/model"
	"github.com/monmaru/gae-graphql/domain/repository"
)

type BlogDatastore struct {
	ProjID string
}

func (b *BlogDatastore) Create(ctx context.Context, blog *model.Blog) (*model.Blog, error) {
	client, err := datastore.NewClient(ctx, b.ProjID)
	if err != nil {
		return nil, err
	}

	key := datastore.IncompleteKey("Blog", nil)
	generatedKey, err := client.Put(ctx, key, blog)
	if err != nil {
		return nil, err
	}

	blog.ID = strconv.FormatInt(generatedKey.ID, 10)
	return blog, nil
}

func (b *BlogDatastore) NewQuery() repository.Query {
	return &QueryImpl{
		projID: b.ProjID,
		query:  datastore.NewQuery("Blog"),
	}
}

type QueryImpl struct {
	projID string
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
	client, err := datastore.NewClient(ctx, q.projID)
	if err != nil {
		return nil, err
	}

	var result model.BlogList
	keys, err := client.GetAll(ctx, q.query, &result.Nodes)
	if err != nil {
		return &result, err
	}

	for i, key := range keys {
		result.Nodes[i].ID = strconv.FormatInt(key.ID, 10)
	}
	result.TotalCount = len(result.Nodes)
	return &result, nil
}
