package database

import (
	"github.com/jinzhu/gorm"
)

// UserInfo 用户实体表
type UserInfo struct {
	gorm.Model
	Account   string  `gorm:"type varchar(10);unique; not null"`
	Type      int8    `gorm:"default:0;not null"`
	Remaining float64 `gorm:"default:0.0; not null"`
	Extra     []byte
}

// CreateUserInfo 创建用户信息表
func CreateUserInfo(userInfo UserInfo) {
	if client.db.NewRecord(userInfo) {
		client.db.Create(&userInfo)
	}
}
