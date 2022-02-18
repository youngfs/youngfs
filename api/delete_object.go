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
	Recursive bool `form:"recursive" json:"recursive" uri:"recursive"`
}

func DeleteObjectHandler(c *gin.Context) {
	mtime := time.Unix(time.Now().Unix(), 0)

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

	set, fp := iam.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
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

	err = directory.DeleteInodeAndEntry(set, fp, mtime, deleteObjectInfo.Recursive)
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
