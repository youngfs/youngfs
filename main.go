package main

import (
	"log"
	"object-storage-server/command"
	"object-storage-server/routers"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	r := routers.InitRouter()
	err := r.Run(":" + command.Port)
	if err != nil {
		log.Println(err)
	}
	return
}
