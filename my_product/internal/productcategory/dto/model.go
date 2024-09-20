package dto

import (
	"time"
)

// 商品分类实体
type ProductCategory struct {
	Id        uint64    `json:"id" xorm:"pk"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (*ProductCategory) TableName() string {
	return "prod_category"
}
