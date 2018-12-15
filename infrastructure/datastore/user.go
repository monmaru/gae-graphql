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
	kind   string
}

func NewUserDatastore(projID string) (*UserDatastore, error) {
	client, err := newDataStoreClient(context.Background(), projID)
	if err != nil {
		return nil, err
	}
	return &UserDatastore{client: client, kind: "User"}, nil
}

func (u *UserDatastore) Create(ctx context.Context, user *model.User) (*model.User, error) {
	key := datastore.IncompleteKey(u.kind, nil)
	generatedKey, err := u.client.Put(ctx, key, user)
	if err != nil {
		return nil, err
	}

	user.ID = strconv.FormatInt(generatedKey.ID, 10)
	return user, nil
}

func (u *UserDatastore) CreateMulti(ctx context.Context, users []*model.User) ([]*model.User, error) {
	var keys []*datastore.Key
	for i := 0; i < len(users); i++ {
		key := datastore.IncompleteKey(u.kind, nil)
		keys = append(keys, key)
	}

	generatedKeys, err := u.client.PutMulti(ctx, keys, users)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(users); i++ {
		users[i].ID = strconv.FormatInt(generatedKeys[i].ID, 10)
	}
	return users, nil
}

func (u *UserDatastore) Get(ctx context.Context, strID string) (*model.User, error) {
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid id")
	}

	user := &model.User{ID: strID}
	key := datastore.IDKey(u.kind, id, nil)
	if err := u.client.Get(ctx, key, user); err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
