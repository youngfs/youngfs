package api

import (
	"github.com/gin-gonic/gin"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"icesos/util"
	"io"
	"net/http"
)

func PutObjectHandler(c *gin.Context) {
	var file io.ReadCloser
	var contentLength uint64
	var filename string
	file, head, err := c.Request.FormFile("file")
	if head != nil && err == nil {
		contentLength = uint64(head.Size)
		filename = head.Filename
	}
	if err != nil {
		file = c.Request.Body
		contentLength = util.GetContentLength(c.Request)
		filename = ""
	}
	defer func() {
		_ = file.Close()
	}()

	setName, fp := set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += full_path.FullPath(filename)
	}
	if !setName.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalSetName]
		c.Set(vars.CodeKey, err.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  err.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}
	if !fp.IsLegalObjectName() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalObjectName]
		c.Set(vars.CodeKey, err.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  err.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	err = server.PutObject(c, setName, fp, contentLength, file)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.Set(vars.CodeKey, err.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  err.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
