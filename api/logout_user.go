package api

import (
	"github.com/gin-gonic/gin"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/iam"
	"net/http"
)

type LogoutUserInfo struct {
	AdminName string `form:"adminName" json:"adminName" uri:"adminName" binding:"required"`
	AdminSK   string `form:"adminSK" json:"adminSK" uri:"adminSK" binding:"required"`
	User      string `form:"user" json:"user" uri:"user" binding:"required"`
}

func LogoutUserHandler(c *gin.Context) {
	logoutUserInfo := &LogoutUserInfo{}

	err := c.Bind(logoutUserInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if logoutUserInfo.AdminName != vars.AdminName || logoutUserInfo.AdminSK != vars.AdminSK {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrAdminAuthenticate].Error(),
			},
		)
		return
	}

	user := iam.User(logoutUserInfo.User)
	_, err = user.Delete()
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
