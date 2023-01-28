package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/full_path"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/log"
	"youngfs/util"
	"youngfs/vars"
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
		apiErr := errors.ErrRouter
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

	set, fp := fs_set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += full_path.FullPath(filename)
	}
	if !set.IsLegal() {
		err := errors.ErrIllegalSetName
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
		err := errors.ErrIllegalObjectName
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

	err = server.PutObject(c, set, fp, contentLength, file, putObjectInfo.Compress)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%+v\n", c.Value(vars.UUIDKey), apiErr)
		}
		c.Set(vars.CodeKey, apiErr.ErrorCode)
		c.Set(vars.ErrorKey, apiErr.Error())
		c.JSON(
			apiErr.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  apiErr.ErrorCode,
				vars.ErrorKey: apiErr.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
