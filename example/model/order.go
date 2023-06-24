// Code generated by gormgen. DO NOT EDIT.
package model

import (
	"time"
)

// Order table model
type Order struct {
	Id         int64     `gorm:"column:id;primary_key" json:"id"`      // 主键
	Name       string    `gorm:"column:name" json:"name"`              // 名称
	Price      float64   `gorm:"column:price" json:"price"`            // 订单价格
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"` // 更新时间
	Deleted    int       `gorm:"column:deleted" json:"deleted"`        // 删除状态，1表示软删除
}

func (order *Order) TableName() string {
	return OrderTableName
}