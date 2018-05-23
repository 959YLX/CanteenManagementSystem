package database

import (
	"github.com/jinzhu/gorm"
)

// GoodsInfo 商品信息实体
type GoodsInfo struct {
	gorm.Model
	Number   int64   `gorm:"AUTO_INCREMENT;unique"`
	BelongTo string  `gorm:"not null"`
	Species  int64   `gorm:"not_null"`
	Price    float64 `gorm:"not null;default:0.0"`
	Extra    []byte
}
