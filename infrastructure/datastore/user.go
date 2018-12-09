package datastore

import (
	"context"
	"errors"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/monmaru/gae-graphql/domain/model"
)

type UserDatastore struct {
	ProjID string
}

func (u *UserDatastore) Create(ctx context.Context, user *model.User) (*model.User, error) {
	client, err := datastore.NewClient(ctx, u.ProjID)
	if err != nil {
		return nil, err
	}

	key := datastore.IncompleteKey("User", nil)
	generatedKey, err := client.Put(ctx, key, user)
	if err != nil {
		return nil, err
	}

	user.ID = strconv.FormatInt(generatedKey.ID, 10)
	return user, nil
}

func (u *UserDatastore) Get(ctx context.Context, strID string) (*model.User, error) {
	client, err := datastore.NewClient(ctx, u.ProjID)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid id")
	}

	user := &model.User{ID: strID}
	key := datastore.IDKey("User", id, nil)
	if err := client.Get(ctx, key, user); err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
