package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/fullpath"
	"youngfs/fs/server"
)

type DeleteObjectInfo struct {
	Recursive bool `form:"recursive" json:"recursive" uri:"recursive" xml:"recursive" yaml:"recursive"`
}

func DeleteObjectHandler(c *gin.Context) {
	deleteObjectInfo := &DeleteObjectInfo{}

	err := c.Bind(deleteObjectInfo)
	if err != nil {
		errorHandler(c, errors.ErrRouter.WrapErrNoStack(err))
		return
	}

	bkt, fp := bucket.Bucket(c.Param("bucket")), fullpath.FullPath(c.Param("path"))
	if !bkt.IsLegal() {
		errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegal() {
		errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	err = server.DeleteObject(c, bkt, fp, deleteObjectInfo.Recursive)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.Status(http.StatusAccepted)
	return
}
