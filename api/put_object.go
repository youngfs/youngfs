package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"icesos/storageEngine"
	"net/http"
	"os"
)

type PutObjectInfo struct {
	BucketName string `form:"BucketName" json:"BucketName" uri:"BucketName" binding:"required"`
	AccessKey  string `form:"AccessKey" json:"AccessKey" uri:"AccessKey" binding:"required"`
	ObjectName string `form:"ObjectName" json:"ObjectName" uri:"ObjectName" binding:"required"`
	DataTime   string `form:"DataTime" json:"DataTime" uri:"DataTime" binding:"required"`
	HostName   string `form:"HostName" json:"HostName" uri:"HostName"`
}

func PutObjectHandler(c *gin.Context) {
	putObjectInfo := &PutObjectInfo{}

	err := c.Bind(putObjectInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			},
		)
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", putObjectInfo)

	assignFileInfo, err := storageEngine.AssignFileHandler()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			},
		)
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", assignFileInfo)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"url":    assignFileInfo.Url,
			"fid":    assignFileInfo.Fid,
		},
	)

	return
}
