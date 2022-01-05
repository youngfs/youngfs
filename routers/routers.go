package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/assignObject", api.AssignObjectHandler)
	r.POST("/putObject", api.PutObjectHandler)
	r.GET("/getObject", api.GetObjectHandler)
	r.GET("/listObject", api.ListObjectHandler)
	r.DELETE("/deleteObject", api.DeleteObjectHandler)
	r.POST("/registerUser", api.RegisterUserHandler)
	r.DELETE("/logoutUser", api.LogoutUserHandler)
	r.POST("/addReadPermission", api.AddReadPermissionHandler)
	r.POST("/addWritePermission", api.AddWritePermissionHandler)
	r.DELETE("/deleteReadPermission", api.DeleteReadPermissionHandler)
	r.DELETE("/deleteWritePermission", api.DeleteWritePermissionHandler)
	return r
}
