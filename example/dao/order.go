package dao

import (
	"context"
	"errors"
	"time"

	"github.com/shinxiang/gormgen/example/model"
	"github.com/shinxiang/gormgen/example/opt"
	"gorm.io/gorm"
)

var _ IOrderDao = (*OrderDao)(nil)

type IOrderDao interface {
	Insert(ctx context.Context, orders ...*model.Order) (err error)
	Save(ctx context.Context, order *model.Order) (err error)
	FindOne(ctx context.Context, condition *opt.OrderOption) (order *model.Order, err error)
	FindList(ctx context.Context, condition *opt.OrderOption) (orders []model.Order, total int64, err error)
	Count(ctx context.Context, condition *opt.OrderOption) (count int64, err error)
	Update(ctx context.Context, order *model.Order, condition *opt.OrderOption) (err error)
	Delete(ctx context.Context, condition *opt.OrderOption) (err error)

	// write you method here
}

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao(db *gorm.DB) *OrderDao {
	return &OrderDao{db: db}
}

func (m *OrderDao) Error(db *gorm.DB) error {
	if db.Error != gorm.ErrRecordNotFound {
		return db.Error
	}
	return nil
}

func (m *OrderDao) Insert(ctx context.Context, orders ...*model.Order) (err error) {
	if orders == nil {
		return errors.New("insert must include orders model")
	}

	for i := range orders {
		if orders[i].CreateTime.IsZero() {
			orders[i].CreateTime = time.Now()
		}
		if orders[i].UpdateTime.IsZero() {
			orders[i].UpdateTime = time.Now()
		}
	}
	db := m.db.WithContext(ctx).Create(orders)
	err = m.Error(db)
	return
}

func (m *OrderDao) Save(ctx context.Context, order *model.Order) (err error) {
	if order == nil {
		return errors.New("save must include order model")
	}

	if order.CreateTime.IsZero() {
		order.CreateTime = time.Now()
	}

	if order.UpdateTime.IsZero() {
		order.UpdateTime = time.Now()
	}

	db := m.db.WithContext(ctx).Save(order)
	err = m.Error(db)
	return
}

func (m *OrderDao) FindOne(ctx context.Context, condition *opt.OrderOption) (order *model.Order, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)

	db = db.First(&order)
	err = m.Error(db)
	return
}

func (m *OrderDao) FindList(ctx context.Context, condition *opt.OrderOption) (orders []model.Order, total int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)

	db = db.Table(model.OrderTableName)
	if err = db.Count(&total).Error; total == 0 {
		return
	}

	if condition != nil {
		db = condition.BuildPage(db)
	}
	db = db.Find(&orders)
	err = m.Error(db)
	return
}

func (m *OrderDao) Count(ctx context.Context, condition *opt.OrderOption) (count int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)

	db = db.Table(model.OrderTableName).Count(&count)
	err = m.Error(db)
	return
}

func (m *OrderDao) Update(ctx context.Context, order *model.Order, condition *opt.OrderOption) (err error) {
	if order == nil {
		return errors.New("update must include order model")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}

	if order.UpdateTime.IsZero() {
		order.UpdateTime = time.Now()
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
	db = db.Where("deleted != ?", 1)

	db = db.Table(model.OrderTableName).Updates(order)
	err = m.Error(db)
	return
}

func (m *OrderDao) Delete(ctx context.Context, condition *opt.OrderOption) (err error) {
	if condition == nil {
		return errors.New("delete must include where condition")
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
	db = db.Where("deleted != ?", 1)

	db = db.Table(model.OrderTableName).
		Select("deleted", "update_time").
		Updates(&model.Order{
			Deleted:    1,
			UpdateTime: time.Now(),
		})

	err = m.Error(db)
	return
}
