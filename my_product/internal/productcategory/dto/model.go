package dto

import (
	"time"
)

// 商品分类实体
type ProductCategory struct {
	Id         uint64    `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"createAt"`
	Updated_at time.Time `json:"updateAt"`
}
