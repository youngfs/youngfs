package api

import (
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/storage_engine"
	"net/http"
)

type AssignObjectInfo struct {
	User       string `form:"user" json:"user" uri:"user" binding:"required"`
	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
	Set        string `form:"set" json:"set" uri:"set" binding:"required"`
	ObjectName string `form:"objectName" json:"objectName" uri:"objectName" binding:"required"`
	FileSize   uint64 `form:"fileSize" json:"fileSize" uri:"fileSize" binding:"required"`
}

func AssignObjectHandler(c *gin.Context) {
	assignObjectInfo := &AssignObjectInfo{}

	err := c.Bind(assignObjectInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	fp := full_path.FullPath(assignObjectInfo.ObjectName)
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

	user := iam.User(assignObjectInfo.User)
	set := iam.Set(assignObjectInfo.Set)
	if !user.Identify(assignObjectInfo.SecretKey) {
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

	assignFileInfo, err := storage_engine.AssignObject(assignObjectInfo.FileSize)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"url": assignFileInfo.Url,
			"fid": assignFileInfo.Fid,
		},
	)
	return
}
