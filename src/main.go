package main

import (
	"os"

	"geekylx.com/CanteenManagementSystemBackend/src/utils"

	"geekylx.com/CanteenManagementSystemBackend/src/service"

	"geekylx.com/CanteenManagementSystemBackend/src/cache"
	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/route"
)

const (
	ROOT_ACCOUNT = "0000000000"
)

func main() {
	args := os.Args
	database.InitDatabase()
	cache.InitRedisClient()
	if len(args) >= 3 && args[1] == "--init" {
		adminPassword := args[2]
		encodedPassword, err := utils.EncodePassword(adminPassword)
		if err != nil || encodedPassword == nil {
			panic(err)
		}
		userInfo := database.UserInfo{
			Account: ROOT_ACCOUNT,
			Type:    uint8(service.USER_TYPE_ROOT),
			Role:    uint64(service.TypeDefaultRole[service.USER_TYPE_ROOT]),
		}
		userLogin := database.UserLogin{
			Account:  ROOT_ACCOUNT,
			Password: *encodedPassword,
		}
		database.CreateUserInfo(userInfo)
		database.CreateUserLogin(userLogin)
	}
	route.StartHTTPServer("127.0.0.1", 9999)
}
