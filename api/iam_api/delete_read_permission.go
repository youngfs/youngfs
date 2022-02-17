package iam_api

import (
	"github.com/gin-gonic/gin"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/iam"
	"net/http"
)

type DeleteReadPermissionInfo struct {
	AdminName string `form:"adminName" json:"adminName" uri:"adminName" binding:"required"`
	AdminSK   string `form:"adminSK" json:"adminSK" uri:"adminSK" binding:"required"`
	User      string `form:"user" json:"user" uri:"user" binding:"required"`
	Set       string `form:"set" json:"set" uri:"set" binding:"required"`
}

func DeleteReadPermissionHandler(c *gin.Context) {
	deleteReadPermissionInfo := &DeleteReadPermissionInfo{}

	err := c.Bind(deleteReadPermissionInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if deleteReadPermissionInfo.AdminName != vars.AdminName || deleteReadPermissionInfo.AdminSK != vars.AdminSK {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrAdminAuthenticate].Error(),
			},
		)
		return
	}

	user := iam.User(deleteReadPermissionInfo.User)
	set := iam.Set(deleteReadPermissionInfo.Set)
	ret, err := user.IsExist()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	if ret == false {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrUserNotExist].Error(),
			},
		)
		return
	}

	err = user.DeleteReadSetPermission(set)
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
