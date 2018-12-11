package datastore

import (
	"context"
	"sync"

	"cloud.google.com/go/datastore"
)

var (
	once     sync.Once
	instance *datastore.Client // singleton instance
)

func newDataStoreClient(ctx context.Context, projID string) (*datastore.Client, error) {
	var err error
	once.Do(func() {
		instance, err = datastore.NewClient(ctx, projID)
	})
	return instance, err
}
