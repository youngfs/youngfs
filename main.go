package main

import (
	"icesos/command/vars"
	"icesos/kv"
	"icesos/routers"
	"log"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	kv.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	r := routers.InitRouter()
	err := r.Run(":" + vars.Port)
	if err != nil {
		log.Println(err)
	}
	return
}
