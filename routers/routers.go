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

	//iam.api
	//r.POST("/registerUser", api.RegisterUserHandler)
	//r.DELETE("/logoutUser", api.LogoutUserHandler)
	//r.POST("/addReadPermission", api.AddReadPermissionHandler)
	//r.POST("/addWritePermission", api.AddWritePermissionHandler)
	//r.DELETE("/deleteReadPermission", api.DeleteReadPermissionHandler)
	//r.DELETE("/deleteWritePermission", api.DeleteWritePermissionHandler)
	return r
}
