package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MYSQL_USER_NAME            = "CanteenManagement"
	MYSQL_PASSWORD             = "123456"
	MYSQL_ADDRESS              = "localhost"
	MYSQL_PORT                 = 3306
	MYSQL_DATABASE             = "CanteenManagementSystem"
	MYSQL_CONNECT_URL_PARTTERN = "%s:%s@%s:%d/%s?charset=utf8&parseTime=True"
)

var client *Client

// Client 数据库连接客户端
type Client struct {
	db *gorm.DB
}

// InitDatabase 初始化数据库连接
func InitDatabase() (err error) {
	if db, err := gorm.Open("mysql", fmt.Sprintf(MYSQL_CONNECT_URL_PARTTERN, MYSQL_USER_NAME, MYSQL_PASSWORD, MYSQL_ADDRESS, MYSQL_PORT, MYSQL_DATABASE)); err == nil {
		client = &Client{
			db: db,
		}
		db.AutoMigrate(&UserInfo{},
			&UserLogin{},
			&GoodsInfo{},
			&FlowingWater{})
	}
	return
}

func (ref *Client) disconnect() error {
	return ref.db.Close()
}
