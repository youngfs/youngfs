package main

import (
	"icesos/log"
	"icesos/routers"
	"icesos/server"
)

func main() {
	log.InitLogger()
	defer log.Sync()

	server.InitServer()

	routers.InitRouter()
	routers.Run()
	return
}
