package routers

import (
	"github.com/gin-gonic/gin"
	"object-storage-server/api"
)

func InitRouter() (r *gin.Engine) {
	r = gin.Default()
	r.GET("/putObject", api.PutObjectHandler)
	r.GET("/getObject", api.GetObjectHandler)
	return
}
