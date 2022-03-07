package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"net/http"
)

func PutObjectHandler(c *gin.Context) {
	ctx := context.Background()

	file, head, err := c.Request.FormFile("file")
	if err != nil {
		apiErr := errors.ErrorCodeResponse[errors.ErrRouter]
		c.JSON(
			apiErr.HTTPStatusCode,
			gin.H{
				"code":  apiErr.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}

	setName, fp := set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += full_path.FullPath(head.Filename)
	}
	if !setName.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalSetName]
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}
	if !fp.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalObjectName]
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	err = server.Svr.PutObject(ctx, setName, fp, uint64(head.Size), file)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
