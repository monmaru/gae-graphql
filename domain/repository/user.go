package repository

import (
	"context"

	"github.com/monmaru/gae-graphql/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	CreateMulti(ctx context.Context, users []*model.User) ([]*model.User, error)
	Get(ctx context.Context, strID string) (*model.User, error)
}
