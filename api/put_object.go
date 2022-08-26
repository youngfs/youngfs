package api

import (
	"github.com/gin-gonic/gin"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/full_path"
	"icesfs/server"
	"icesfs/set"
	"icesfs/util"
	"io"
	"net/http"
)

type PutObjectInfo struct {
	Recover  bool `form:"recover" json:"recover" uri:"recover"`
	Compress bool `form:"compress" json:"compress" uri:"compress"`
}

func PutObjectHandler(c *gin.Context) {
	putObjectInfo := &PutObjectInfo{
		Recover: false,
	}

	err := c.Bind(putObjectInfo)
	if err != nil {
		apiErr := errors.GetAPIErr(errors.ErrRouter)
		c.Set(vars.CodeKey, apiErr.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			apiErr.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  apiErr.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}

	if putObjectInfo.Recover {
		RecoverObjectHandler(c)
		return
	}

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
		contentLength = util.GetContentLength(c.Request.Header)
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
		err := errors.GetAPIErr(errors.ErrIllegalSetName)
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
		err := errors.GetAPIErr(errors.ErrIllegalObjectName)
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

	err = server.PutObject(c, setName, fp, contentLength, file, putObjectInfo.Compress)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.GetAPIErr(errors.ErrServer)
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
