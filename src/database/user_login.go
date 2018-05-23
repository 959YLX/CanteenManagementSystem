package database

import (
	"github.com/jinzhu/gorm"
)

// UserLogin 用户登录权限关系表
type UserLogin struct {
	gorm.Model
	Account  string `gorm:"type varchar(10); unique;not null"`
	Password string `gorm:"not null"`
	Role     int8   `gorm:"not null;default:0"`
}

// CreateUserLogin 创建用户登录关系表
func CreateUserLogin(userLogin UserLogin) {
	if client.db.NewRecord(userLogin) {
		client.db.Create(&userLogin)
	}
}
