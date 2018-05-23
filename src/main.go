package main

import (
	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/route"
)

func main() {
	database.InitDatabase()
	route.StartHTTPServer("127.0.0.1", 9999)
}
