package database

import (
	"github.com/jinzhu/gorm"
)

// UserLogin 用户登录权限关系表
type UserLogin struct {
	gorm.Model
	Account  string `gorm:"type varchar(10); unique;not null"`
	Password string `gorm:"not null"`
	Role     uint8  `gorm:"not null;default:0"`
}

// CreateUserLogin 创建用户登录关系表
func CreateUserLogin(userLogin UserLogin) {
	if client.db.NewRecord(userLogin) {
		client.db.Create(&userLogin)
	}
}

// GetUserLoginByAccount 通过account获取用户登录关系表
func GetUserLoginByAccount(account string) (*UserLogin, error) {
	userLogin := &UserLogin{}
	r := client.db.Where("account = ?", account).First(userLogin)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.Error
	}
	return userLogin, nil
}

// GetMaxAccount 获取最大account
func GetMaxAccount() (*string, error) {
	userLogin := &UserLogin{}
	r := client.db.Order("account desc").First(userLogin)
	if r.Error != nil {
		return nil, nil
	}
	account := userLogin.Account
	return &account, nil
}
