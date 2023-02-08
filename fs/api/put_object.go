package api

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/full_path"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/log"
	"youngfs/vars"
)

type PutObjectInfo struct {
	Compress bool `form:"compress" json:"compress" uri:"compress"`
}

func PutObjectHandler(c *gin.Context) {
	putObjectInfo := &PutObjectInfo{}

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

	var file io.ReadCloser
	var filename string
	file, head, err := c.Request.FormFile("file")
	if head != nil && err == nil {
		filename = head.Filename
		if head.Header.Get("Content-Encoding") == "gzip" {
			file, err = gzip.NewReader(file)
			if err != nil {
				err := errors.ErrContentEncoding
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
		}
	}
	if err != nil {
		file = c.Request.Body
		filename = ""
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			file, err = gzip.NewReader(file)
			if err != nil {
				err := errors.ErrContentEncoding
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
		}
	}
	defer func() {
		_ = file.Close()
	}()

	for _, transfer := range c.Request.TransferEncoding {
		if transfer == "chunked" {
			err := errors.ErrNotSupportChunkTransfer
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
	}

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

	err = server.PutObject(c, set, fp, file)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%v\n", c.Value(vars.UUIDKey), err)
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
