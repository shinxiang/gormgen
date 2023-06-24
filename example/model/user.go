// Code generated by gormgen. DO NOT EDIT.
package model

import (
	"time"
)

// User table model
type User struct {
	Id         int64     `gorm:"column:id;primary_key" json:"id"`      // 用户id
	Nickname   string    `gorm:"column:nickname" json:"nickname"`      // 昵称
	Username   string    `gorm:"column:username" json:"username"`      // 用户名
	Password   string    `gorm:"column:password" json:"password"`      // 登录密码
	Salt       string    `gorm:"column:salt" json:"salt"`              // 随机盐
	Phone      string    `gorm:"column:phone" json:"phone"`            // 手机号
	Status     string    `gorm:"column:status" json:"status"`          // 状态：Y-启用；N-禁用
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
	Deleted    int       `gorm:"column:deleted" json:"deleted"`        // 逻辑删除：0-正常，1-删除
}

func (user *User) TableName() string {
	return UserTableName
}