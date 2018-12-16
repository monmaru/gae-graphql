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

func NewUserDatastore(client *datastore.Client) *UserDatastore {
	return &UserDatastore{client: client, kind: "User"}
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

func (u *UserDatastore) GetMulti(ctx context.Context, strIDs []string) ([]*model.User, error) {
	var keys []*datastore.Key
	for _, strID := range strIDs {
		id, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			return nil, errors.New("Invalid id")
		}
		key := datastore.IDKey(u.kind, id, nil)
		keys = append(keys, key)
	}

	var temp = make([]*model.User, len(keys))
	if err := u.client.GetMulti(ctx, keys, temp); err != nil {
		return nil, errors.New("User not found")
	}

	var users []*model.User
	for i := 0; i < len(temp); i++ {
		if temp[i] == nil {
			continue
		}
		temp[i].ID = strIDs[i]
		users = append(users, temp[i])
	}
	return users, nil
}
