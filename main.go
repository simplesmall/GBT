package main

import (
	"GBT/Config"
	"GBT/routers"
)

func main() {
	defer Config.CloseDB()
	routers.InitServer()
}
