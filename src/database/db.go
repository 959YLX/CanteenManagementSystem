package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MYSQL_USER_NAME            = "CanteenManagement"
	MYSQL_PASSWORD             = "123456"
	MYSQL_DATABASE             = "CanteenManagementSystem"
	MYSQL_CONNECT_URL_PARTTERN = "%s:%s@/%s?charset=utf8&parseTime=True"
)

var client *Client

// Client 数据库连接客户端
type Client struct {
	db *gorm.DB
}

// InitDatabase 初始化数据库连接
func InitDatabase() (err error) {
	db, err := gorm.Open("mysql", fmt.Sprintf(MYSQL_CONNECT_URL_PARTTERN, MYSQL_USER_NAME, MYSQL_PASSWORD, MYSQL_DATABASE))
	if err != nil {
		panic(err)
	}
	client = &Client{
		db: db.Debug(),
	}
	db.AutoMigrate(&UserInfo{},
		&UserLogin{},
		&GoodsInfo{},
		&FlowingWater{})
	return
}

// Disconnect 关闭数据库连接
func Disconnect() error {
	return client.db.Close()
}

// Transaction 获取Transaction
func Transaction() *gorm.DB {
	return client.db.Begin()
}
