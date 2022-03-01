package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/:set/*fp", api.PutObjectHandler)
	r.GET("/:set/*fp", api.GetObjectHandler)
	r.DELETE("/:set/*fp", api.DeleteObjectHandler)
	r.HEAD("/:set/*fp", api.HeadObjectHandler)
	return r
}
