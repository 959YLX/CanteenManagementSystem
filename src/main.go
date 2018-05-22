package main

import (
	"geekylx.com/CanteenManagementSystemBackend/src/route"
)

func main() {
	route.StartHTTPServer("127.0.0.1", 9999)
}
