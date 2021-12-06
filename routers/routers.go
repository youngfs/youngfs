package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/putObject", api.PutObjectHandler)
	r.GET("/getObject", api.GetObjectHandler)
	return r
}
