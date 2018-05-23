package database

import (
	"github.com/jinzhu/gorm"
)

// FlowingWater 流水记录实体
type FlowingWater struct {
	gorm.Model
	From    string `gorm:"type varchar(10);unique_index:unique_from_to"`
	To      string `gorm:"type varchar(10);unique_index:unique_from_to"`
	Type    int8   `gorm:"not null"`
	Species int64
	Extra   []byte
}
