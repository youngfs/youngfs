package api

import (
	"github.com/gin-gonic/gin"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/iam"
	"net/http"
)

type RegisterUserInfo struct {
	AdminName string `form:"adminName" json:"adminName" uri:"adminName" binding:"required"`
	AdminSK   string `form:"adminSK" json:"adminSK" uri:"adminSK" binding:"required"`
	User      string `form:"user" json:"user" uri:"user" binding:"required"`
	SecretKey string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
}

func RegisterUserHandler(c *gin.Context) {
	registerUserInfo := &RegisterUserInfo{}

	err := c.Bind(registerUserInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if registerUserInfo.AdminName != vars.AdminName || registerUserInfo.AdminSK != vars.AdminSK {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrAdminAuthenticate].Error(),
			},
		)
		return
	}

	user := iam.User(registerUserInfo.User)
	err = user.Create(registerUserInfo.SecretKey)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
