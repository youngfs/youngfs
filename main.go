package main

import (
	"icesos/command/vars"
	"icesos/kv/redis"
	"icesos/routers"
	"log"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	redis.Client.Initialize(vars.RedisHostPost, vars.RedisPassword, vars.RedisDatabase)
	r := routers.InitRouter()
	err := r.Run(":" + vars.Port)
	if err != nil {
		log.Println(err)
	}
	return
}
