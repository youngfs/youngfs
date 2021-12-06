package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type GetObjectInfo struct {
	BucketName string `form:"BucketName" json:"BucketName" uri:"BucketName" binding:"required"`
	AccessKey  string `form:"AccessKey" json:"AccessKey" uri:"AccessKey" binding:"required"`
	ObjectName string `form:"ObjectName" json:"ObjectName" uri:"ObjectName" binding:"required"`
	DataTime   string `form:"DataTime" json:"DataTime" uri:"DataTime" binding:"required"`
}

func GetObjectHandler(c *gin.Context) {
	getObjectInfo := &GetObjectInfo{}

	err := c.Bind(getObjectInfo)
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

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", getObjectInfo)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
		},
	)

	return
}
