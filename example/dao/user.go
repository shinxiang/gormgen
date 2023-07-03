package dao

import (
	"context"
	"errors"
	"time"

	"github.com/shinxiang/gormgen/example/model"
	"github.com/shinxiang/gormgen/example/opt"
	"gorm.io/gorm"
)

var _ IUserDao = (*UserDao)(nil)

type IUserDao interface {
	Insert(ctx context.Context, users ...*model.User) (err error)
	Save(ctx context.Context, user *model.User) (err error)
	FindOne(ctx context.Context, condition *opt.UserOption) (user *model.User, err error)
	FindList(ctx context.Context, condition *opt.UserOption) (users []model.User, total int64, err error)
	Count(ctx context.Context, condition *opt.UserOption) (count int64, err error)
	Update(ctx context.Context, user *model.User, condition *opt.UserOption) (err error)
	Delete(ctx context.Context, condition *opt.UserOption) (err error)

	// write you method here
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (m *UserDao) Error(db *gorm.DB) error {
	if db.Error != gorm.ErrRecordNotFound {
		return db.Error
	}
	return nil
}

// TableName Get table name from context, key is TABLE_NAME.
func (m *UserDao) TableName(ctx context.Context) string {
	if ctx != nil {
		val := ctx.Value("TABLE_NAME")
		switch val.(type) {
		case string:
			if tableName := val.(string); tableName != "" {
				return tableName
			}
		}
	}
	return model.UserTableName
}

func (m *UserDao) Insert(ctx context.Context, users ...*model.User) (err error) {
	if users == nil {
		return errors.New("insert must include users model")
	}

	for i := range users {
		if users[i].CreateTime.IsZero() {
			users[i].CreateTime = time.Now()
		}
		if users[i].UpdateTime.IsZero() {
			users[i].UpdateTime = time.Now()
		}
	}
	db := m.db.WithContext(ctx).Table(m.TableName(ctx)).Create(users)
	err = m.Error(db)
	return
}

func (m *UserDao) Save(ctx context.Context, user *model.User) (err error) {
	if user == nil {
		return errors.New("save must include user model")
	}

	if user.CreateTime.IsZero() {
		user.CreateTime = time.Now()
	}

	if user.UpdateTime.IsZero() {
		user.UpdateTime = time.Now()
	}

	db := m.db.WithContext(ctx).Table(m.TableName(ctx)).Save(user)
	err = m.Error(db)
	return
}

func (m *UserDao) FindOne(ctx context.Context, condition *opt.UserOption) (user *model.User, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)

	db = db.Table(m.TableName(ctx)).First(&user)
	err = m.Error(db)
	return
}

func (m *UserDao) FindList(ctx context.Context, condition *opt.UserOption) (users []model.User, total int64, err error) {
	db := m.db.WithContext(ctx).Table(m.TableName(ctx))
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)
	if err = db.Count(&total).Error; total == 0 {
		return
	}

	if condition != nil {
		db = condition.BuildPage(db)
	}
	db = db.Find(&users)
	err = m.Error(db)
	return
}

func (m *UserDao) Count(ctx context.Context, condition *opt.UserOption) (count int64, err error) {
	db := m.db.WithContext(ctx)
	if condition != nil {
		db = condition.BuildQuery(db)
	}
	db = db.Where("deleted != ?", 1)

	db = db.Table(m.TableName(ctx)).Count(&count)
	err = m.Error(db)
	return
}

func (m *UserDao) Update(ctx context.Context, user *model.User, condition *opt.UserOption) (err error) {
	if user == nil {
		return errors.New("update must include user model")
	} else if condition == nil {
		return errors.New("update must include where condition")
	}

	if user.UpdateTime.IsZero() {
		user.UpdateTime = time.Now()
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
	db = db.Where("deleted != ?", 1)

	db = db.Table(m.TableName(ctx)).Updates(user)
	err = m.Error(db)
	return
}

func (m *UserDao) Delete(ctx context.Context, condition *opt.UserOption) (err error) {
	if condition == nil {
		return errors.New("delete must include where condition")
	}

	db := m.db.WithContext(ctx)
	db = condition.BuildQuery(db)
	db = db.Where("deleted != ?", 1)

	db = db.Table(m.TableName(ctx)).
		Select("deleted", "update_time").
		Updates(&model.User{
			Deleted:    1,
			UpdateTime: time.Now(),
		})

	err = m.Error(db)
	return
}
