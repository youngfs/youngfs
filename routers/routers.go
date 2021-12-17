package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api/get_object"
	"icesos/api/put_object"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/putObject", put_object.PutObjectHandler)
	r.GET("/getObject", get_object.GetObjectHandler)
	return r
}
