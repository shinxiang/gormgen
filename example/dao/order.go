package dao

import (
	"context"

	"github.com/shinxiang/gormgen/example/model"
)

var _ IOrderDao = (*OrderDao)(nil)

type IOrderDao interface {
	Insert(ctx context.Context, orders ...*model.Order) (err error)
	Save(ctx context.Context, order *model.Order) (err error)
	FindOne(ctx context.Context, condition *model.OrderOption) (order *model.Order, err error)
	FindList(ctx context.Context, condition *model.OrderOption) (orders []model.Order, total int64, err error)
	Count(ctx context.Context, condition *model.OrderOption) (count int64, err error)
	Update(ctx context.Context, order *model.Order, condition *model.OrderOption) (err error)
	Delete(ctx context.Context, condition *model.OrderOption) (err error)

	// write you method here
}
