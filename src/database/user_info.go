package database

import (
	"github.com/jinzhu/gorm"
)

// UserInfo 用户实体表
type UserInfo struct {
	gorm.Model
	Account   string  `gorm:"type varchar(10);unique; not null"`
	Type      uint8   `gorm:"default:0;not null"`
	Remaining float64 `gorm:"default:0.0; not null"`
	Extra     []byte
}

// CreateUserInfo 创建用户信息表
func CreateUserInfo(userInfo UserInfo) {
	if client.db.NewRecord(userInfo) {
		client.db.Create(&userInfo)
	}
}

// GetUserInfoByAccount 通过账号获取用户信息
func GetUserInfoByAccount(account string) (*UserInfo, error) {
	userInfo := &UserInfo{}
	r := client.db.Where("account = ?", account).First(userInfo)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.Error
	}
	return userInfo, nil
}
