package main

import (
	"icesos/command"
	"icesos/routers"
	"log"
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
