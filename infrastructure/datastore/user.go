package datastore

import (
	"context"
	"errors"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/monmaru/gae-graphql/domain/model"
)

type UserDatastore struct {
	client *datastore.Client
}

func NewUserDatastore(projID string) (*UserDatastore, error) {
	client, err := newDataStoreClient(context.Background(), projID)
	if err != nil {
		return nil, err
	}
	return &UserDatastore{client: client}, nil
}

func (u *UserDatastore) Create(ctx context.Context, user *model.User) (*model.User, error) {
	key := datastore.IncompleteKey("User", nil)
	generatedKey, err := u.client.Put(ctx, key, user)
	if err != nil {
		return nil, err
	}

	user.ID = strconv.FormatInt(generatedKey.ID, 10)
	return user, nil
}

func (u *UserDatastore) Get(ctx context.Context, strID string) (*model.User, error) {
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
