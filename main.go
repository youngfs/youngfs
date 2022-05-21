package main

import (
	"icesos/command"
	"icesos/log"
	"icesos/routers"
	"icesos/server"
)

func main() {
	command.InitCommand()

	log.InitLogger()
	defer log.Sync()

	server.InitServer()

	routers.InitRouter()
	routers.Run()
	return
}
