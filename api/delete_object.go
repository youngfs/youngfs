package api

import (
	"github.com/gin-gonic/gin"
	"icesos/directory"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"net/http"
	"time"
)

type DeleteObjectInfo struct {
	User       string `form:"user" json:"user" uri:"user" binding:"required"`
	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
	Set        string `form:"set" json:"set" uri:"set" binding:"required"`
	ObjectName string `form:"objectName" json:"objectName" uri:"objectName" binding:"required"`
	Recursive  bool   `form:"recursive" json:"recursive" uri:"recursive"`
}

func DeleteObjectHandler(c *gin.Context) {
	deleteObjectInfo := &DeleteObjectInfo{}

	err := c.Bind(deleteObjectInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	set := iam.Set(deleteObjectInfo.Set)
	fp := full_path.FullPath(deleteObjectInfo.ObjectName)
	if !fp.IsLegal() {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrIllegalObjectName].Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	user := iam.User(deleteObjectInfo.User)
	if !user.Identify(deleteObjectInfo.SecretKey) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrUserAuthenticate].Error(),
			},
		)
		return
	}

	if !user.WriteSetPermission(set) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrSetWriteAuthenticate].Error(),
			},
		)
		return
	}

	err = directory.DeleteInodeAndEntry(set, fp, time.Unix(time.Now().Unix(), 0), deleteObjectInfo.Recursive)
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
