package main

import (
	"icesfs/command"
	"icesfs/log"
	"icesfs/routers"
	"icesfs/server"
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
