package database

import (
	"github.com/jinzhu/gorm"
)

// UserInfo 用户实体表
type UserInfo struct {
	gorm.Model
	Account   string  `gorm:"type varchar(10);unique; not null"`
	Type      uint8   `gorm:"default:0;not null"`
	Role      uint64  `gorm:"not null;default:0"`
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

// GetUserInfosByAccounts 通过account数组获取UserInfo数组
func GetUserInfosByAccounts(accounts []string) (ref []*UserInfo, err error) {
	r := client.db.Where("account in (?)", accounts).Find(&ref)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.Error
	}
	return
}

// DeleteUserInfosByAccountsInTransaction 软删除用户信息(事务)
func DeleteUserInfosByAccountsInTransaction(tx *gorm.DB, accounts []string) error {
	r := tx.Where("account in (?)", accounts).Delete(&UserInfo{})
	if r.Error != nil {
		tx.Rollback()
		return r.Error
	}
	return nil
}
