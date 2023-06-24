package dao

import (
	"context"

	"github.com/shinxiang/gormgen/example/model"
)

var _ IUserDao = (*UserDao)(nil)

type IUserDao interface {
	Insert(ctx context.Context, users ...*model.User) (err error)
	Save(ctx context.Context, user *model.User) (err error)
	FindOne(ctx context.Context, condition *model.UserOption) (user *model.User, err error)
	FindList(ctx context.Context, condition *model.UserOption) (users []model.User, total int64, err error)
	Count(ctx context.Context, condition *model.UserOption) (count int64, err error)
	Update(ctx context.Context, user *model.User, condition *model.UserOption) (err error)
	Delete(ctx context.Context, condition *model.UserOption) (err error)

	// write you method here
}
