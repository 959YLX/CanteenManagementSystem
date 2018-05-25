package database

import (
	"github.com/jinzhu/gorm"
)

// FlowingWater 流水记录实体
type FlowingWater struct {
	gorm.Model
	From    string  `gorm:"type varchar(10);index:index_from_to"`
	To      string  `gorm:"type varchar(10);index:index_from_to"`
	Type    uint8   `gorm:"not null"`
	Money   float64 `gorm:"not null"`
	Species uint64
	Extra   []byte
}

// RecordInTransaction 记录交易流水(事务)
func RecordInTransaction(tx *gorm.DB, flowingWater FlowingWater) error {
	r := tx.Create(&flowingWater)
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	return nil
}
