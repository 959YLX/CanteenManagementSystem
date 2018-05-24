package database

import (
	"github.com/jinzhu/gorm"
)

// GoodsInfo 商品信息实体
type GoodsInfo struct {
	gorm.Model
	Number   uint64  `gorm:"AUTO_INCREMENT;unique"`
	BelongTo string  `gorm:"not null"`
	Species  uint64  `gorm:"not_null"`
	Price    float64 `gorm:"not null;default:0.0"`
	Extra    []byte
}

// SaveGoods 保存商品信息
func SaveGoods(goods GoodsInfo) error {
	r := client.db.Create(goods)
	return r.Error
}
