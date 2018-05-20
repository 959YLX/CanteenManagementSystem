package main

import "./route"

func main() {
	route.StartHttpServer("127.0.0.1", 9999)
}
